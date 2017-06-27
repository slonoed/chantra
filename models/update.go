package models

import (
	"errors"
	"log"

	"github.com/slonoed/chantra/state"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	tg "gopkg.in/telegram-bot-api.v4"
)

// Update represents telegram update + id in mongo
type Update struct {
	tg.Update
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BotID bson.ObjectId `bson:"botId" json:"-"`
}

func CreateUpdate(tu tg.Update, bot *Bot) Update {
	return Update{tu, "", bot.ID}
}

// Save update to database
func (update *Update) Save(session *mgo.Session) error {
	_, err := session.DB(dbName).C("updates").Upsert(bson.M{"update_id": update.UpdateID}, update)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (update *Update) IsAlreadyVoted(app *state.AppState) (bool, error) {
	if update.CallbackQuery == nil {
		return true, errors.New("Not callbackquery")
	}
	if update.CallbackQuery.Message == nil {
		return true, errors.New("No message")
	}
	session := app.CloneSession()
	defer session.Close()
	mid := update.CallbackQuery.Message.MessageID
	uid := update.CallbackQuery.From.ID
	query := bson.M{
		"update.callbackquery.message.messageid": mid,
		"update.callbackquery.from.id":           uid,
	}
	count, err := session.DB(dbName).C("updates").Find(query).Count()
	if err != nil {
		return true, err
	}
	notVoted := count == 1 // TODO check id?
	return !notVoted, nil
}
