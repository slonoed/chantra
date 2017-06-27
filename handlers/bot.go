package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/slonoed/chntr/models"
	"bitbucket.org/slonoed/chntr/state"
)

type botInput struct {
	Token string `json:"token"`
}
type botOut struct {
	OK  bool        `json:"ok"`
	Bot *models.Bot `json:"bot"`
}

// HandleAddBot create bot
func AddBot(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant read body"}`)
		return
	}

	var input botInput
	if err = json.Unmarshal(data, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Body should be valid json"}`)
		return
	}

	user := r.Context().Value(AuthContextKey("user")).(*models.User)

	bot, err := models.GetBotByToken(app, input.Token, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Wrong token"}`)
		return
	}
	out := &botOut{true, bot}
	bytes, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant marshal bot"}`)
		return
	}

	fmt.Fprint(w, string(bytes))
}

type listBotsOut struct {
	OK   bool         `json:"ok"`
	Bots []models.Bot `json:"bots"`
}

// ListBots create bot
func ListBots(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(AuthContextKey("user")).(*models.User)
	bots, err := models.GetBotsForUser(app.MgoSession, user)

	out := &listBotsOut{true, bots}
	bytes, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant marshal bot"}`)
		return
	}
	fmt.Fprint(w, string(bytes))
}
