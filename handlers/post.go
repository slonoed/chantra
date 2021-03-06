package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/slonoed/chantra/models"
	"github.com/slonoed/chantra/state"
)

type addPostOut struct {
	OK   bool         `json:"ok"`
	Post *models.Post `json:"post"`
}

// AddPost and schedule sending
func AddPost(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant read body"}`)
		return
	}

	var post models.Post
	if err = json.Unmarshal(data, &post); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Body should be valid json"}`)
		return
	}

	if post.Text == "" {
		fmt.Fprintf(w, `{"ok": false, "error": "Empty text"}`)
		return
	}

	fmt.Println("Whut", post.ChannelID)
	channel, err := models.GetChannelByID(app, post.ChannelID)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Server error"}`)
		return
	}

	bot, err := models.GetBotById(app, channel.BotID.Hex())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"ok": false, "error": "No bot"}`)
		return
	}
	user := r.Context().Value(AuthContextKey("user")).(*models.User)
	if bot.UserID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"ok": false, "error": "Not authorized to use this bot"}`)
		return
	}
	err = post.Save(app)
	if err != nil {
		fmt.Fprintf(w, `{"ok": false, "error": "Server error"}`)
		return
	}

	out := &addPostOut{true, &post}
	bytes, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant marshal post"}`)
		return
	}
	fmt.Fprint(w, string(bytes))
}
