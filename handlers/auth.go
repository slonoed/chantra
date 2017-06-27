package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.org/slonoed/chntr/models"
	"bitbucket.org/slonoed/chntr/state"
	"golang.org/x/oauth2"
)

type AuthContextKey string

// FIXME
var oauthStateString = "rand"

func GoogleLogin(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	url := app.OauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := app.OauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	// get from google profile
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	profile := models.GoogleProfile{}
	if err = json.Unmarshal(data, &profile); err != nil {
		log.Println("Unm", err)
		// FIXME add response
		return
	}

	// get from or add to database
	user, err := models.UserFromProfile(app.MgoSession, profile)
	if err != nil {
		log.Println("Cant add", err)
		// FIXME
		return
	}

	cookieSession, err := app.SessionStore.Get(r, "chantra")
	if err != nil {
		log.Println("What", err)
		// FIXME
		return
	}

	fmt.Println("Set cookie ", user)
	cookieSession.Values["id"] = user.ID.Hex()
	cookieSession.Options.HttpOnly = true
	err = cookieSession.Save(r, w)
	if err != nil {
		log.Println("Cookie get", err)
		// FIXME
		return
	}
	fmt.Println("no5")

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func Logout(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	cookieSession, err := app.SessionStore.Get(r, "chantra")
	if err != nil {
		log.Println("handleLogout store.Get", err)
		// FIXME
		return
	}
	cookieSession.Values["id"] = ""
	cookieSession.Options.HttpOnly = true
	err = cookieSession.Save(r, w)
	if err != nil {
		log.Println("Cookie get", err)
		// FIXME
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
