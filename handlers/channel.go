package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slonoed/chantra/models"
	"github.com/slonoed/chantra/state"
	tg "gopkg.in/telegram-bot-api.v4"
)

type channelInput struct {
	Name  string `json:"name"`
	BotID string `json:"botId"`
}

type channelOut struct {
	OK      bool            `json:"ok"`
	Channel *models.Channel `json:"channel"`
}

func AddChannel(app *state.AppState, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant read body"}`)
		return
	}

	var input channelInput
	if err = json.Unmarshal(data, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "Body should be valid json"}`)
		return
	}

	fmt.Println("WUT")
	user := r.Context().Value(AuthContextKey("user")).(*models.User)
	bot, err := models.GetBotById(app, input.BotID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Some error occurred"}`)
		return
	}
	if bot.UserID != user.ID {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok": false, "error": "You not an owner of this bot"}`)
		return
	}

	botAPI, err := tg.NewBotAPI(bot.Token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Some error occurred"}`)
		return
	}

	chat, err := botAPI.GetChat(tg.ChatConfig{0, input.Name})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Some error occurred"}`)
		return
	}

	channel, err := models.CreateChannel(app, chat, bot.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Some error occurred"}`)
		return
	}

	out := &channelOut{true, channel}
	bytes, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"ok": false, "error": "Cant marshal bot"}`)
		return
	}

	fmt.Fprint(w, string(bytes))
}

/*
{
  "ok": true,
  "result": {
    "id": -1001105623419,
    "title": "My test channel",
    "type": "channel",
    "username": "slonoed_test_channel"
  }
}
*/
