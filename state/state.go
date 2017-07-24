package state

import (
	"html/template"
	"net/http"

	"github.com/kidstuff/mongostore"
	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
)

const EnvDev = "development"
const EnvProd = "production"

// AppState contains all current state of app
type AppState struct {
	MgoSession   *mgo.Session
	SessionStore *mongostore.MongoStore
	Templates    map[string]*template.Template
	OauthConfig  *oauth2.Config
	Done         chan struct{}
	Env          string
}

// AppHandler incapsulates http handler + app state
type AppHandler struct {
	*AppState
	H func(*AppState, http.ResponseWriter, *http.Request)
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO add try catch
	ah.H(ah.AppState, w, r)
}

func (app *AppState) CloneSession() *mgo.Session {
	return app.MgoSession.Clone()
}
