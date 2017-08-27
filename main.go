package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kidstuff/mongostore"
	h "github.com/slonoed/chantra/handlers"
	"github.com/slonoed/chantra/models"
	"github.com/slonoed/chantra/state"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	ms, err := mgo.Dial(dbURL)
	if err != nil {
		log.Panicf("Dial mongo error: %s. URL: %s", err, dbURL)
	}
	defer ms.Close()

	sessionColl := ms.DB("chantra").C("sessions")
	store := mongostore.NewMongoStore(sessionColl, 3600*24*7, true, []byte(os.Getenv("SESSION_SECRET")))

	t, _ := template.ParseFiles("views/index.html")
	if err != nil {
		log.Panicln(err)
	}
	indexTemplate := t
	templates := map[string]*template.Template{"index": indexTemplate}
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("G_ADDRESS") + "/auth/google/callback",
		ClientID:     os.Getenv("G_KEY"),
		ClientSecret: os.Getenv("G_SECRET"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	env := state.EnvProd
	if os.Getenv("GO_ENV") == "development" {
		env = state.EnvDev
	}

	appState := &state.AppState{
		MgoSession:   ms,
		SessionStore: store,
		Templates:    templates,
		OauthConfig:  googleOauthConfig,
		Done:         make(chan struct{}),
		Env:          env,
	}

	r := mux.NewRouter()

	r.HandleFunc("/status", handleStatus)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/auth/google", state.AppHandler{appState, h.GoogleLogin})
	r.Handle("/auth/google/callback", state.AppHandler{appState, h.GoogleCallback})
	r.Handle("/auth/logout", state.AppHandler{appState, h.Logout})

	// API
	r.Handle("/api/addBot", apiMiddleware(appState, state.AppHandler{appState, h.AddBot}))
	r.Handle("/api/listBots", apiMiddleware(appState, state.AppHandler{appState, h.ListBots}))
	r.Handle("/api/addPost", apiMiddleware(appState, state.AppHandler{appState, h.AddPost}))
	r.Handle("/api/addChannel", apiMiddleware(appState, state.AppHandler{appState, h.AddChannel}))

	r.PathPrefix("/").Handler(state.AppHandler{appState, handleIndex})

	bots, err := models.GetAllBots(ms)
	if err == nil {
		for _, bot := range bots {
			go func(bot models.Bot) {
				bot.CollectUpdates(appState)
			}(bot)
		}
	}

	go cronPosts(appState)

	// TODO add handlers.RecoveryHandler()(r)
	go http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
	fmt.Println("Server started")

	// Shut down when SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	close(appState.Done)
	fmt.Print("\nExiting")
	for i := 1; i < 4; i++ {
		fmt.Print("\rExiting" + repeatSym(".", i))
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println()
}

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func handleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func apiMiddleware(app *state.AppState, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.User
		cookieSession, err := app.SessionStore.Get(r, "chantra")
		if err == nil {
			sessionID := cookieSession.Values["id"]
			if sessionID != nil {
				ID := sessionID.(string)
				user, err = models.UserByID(app.MgoSession, ID)
			}
		}
		if err != nil || user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"ok": false, "error": "Unauthorized"}`)
			return
		}

		ctx := context.WithValue(r.Context(), h.AuthContextKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type initialData struct {
	User     *models.User     `json:"user"`
	Bots     []models.Bot     `json:"bots"`
	Channels []models.Channel `json:"channels"`
	Posts    []models.Post    `json:"posts"`
}

// PageData struct for rendering page
type PageData struct {
	User *models.User `json:"user"`
	JSON string       `json:"json"`
	Env  string       `json:"env"`
}

func handleIndex(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	var user *models.User
	cookieSession, err := app.SessionStore.Get(r, "chantra")
	if err == nil {
		sessionID := cookieSession.Values["id"]
		if sessionID != nil {
			ID := sessionID.(string)
			user, err = models.UserByID(app.MgoSession, ID)
		}
	}
	var bots []models.Bot
	if user != nil {
		bots, err = models.GetBotsForUser(app.MgoSession, user)
	}
	if err != nil {
		log.Panicln(err)
	}
	channels, err := models.GetChannelsForBots(app, bots)
	if err != nil {
		channels = []models.Channel{}
	}

	posts, err := models.GetPostsForChannels(app, channels)
	if err != nil {
		posts = []models.Post{}
	}

	initialD := initialData{user, bots, channels, posts}
	initial, err := json.Marshal(initialD)
	if err != nil {
		fmt.Println(err)
	}

	data := PageData{user, string(initial), app.Env}
	t, _ := template.ParseFiles("views/index.html")
	t.Execute(w, data)
}

func repeatSym(s string, num int) string {
	a := ""
	for i := 0; i < num; i++ {
		a = a + s
	}
	return a
}

func cronPosts(app *state.AppState) {
	for {
		time.Sleep(time.Second)
		posts, err := models.GetReadyToSendPosts(app)
		if err != nil {
			continue
		}
		for _, post := range posts {
			go func(p models.Post) {
				if err := p.Send(app); err != nil {
					log.Println(err)
				}
			}(post)
		}
	}
}
