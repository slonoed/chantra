package models

import (
	"fmt"
	"log"

	"github.com/slonoed/chantra/state"
	"gopkg.in/mgo.v2/bson"
	tg "gopkg.in/telegram-bot-api.v4"
)

// Channel represents channel chat in mongo
type Channel struct {
	tg.Chat
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BotID bson.ObjectId `bson:"botId" json:"botId"`
}

// CreateChannel for bot
func CreateChannel(app *state.AppState, chat tg.Chat, botID bson.ObjectId) (*Channel, error) {
	channel := &Channel{
		Chat:  chat,
		BotID: botID,
	}
	session := app.CloneSession()
	defer session.Close()

	// Avoid duplicates
	count, _ := session.DB(dbName).C("channels").Find(bson.M{"chat.id": chat.ID}).Count()
	if count != 0 {
		return nil, fmt.Errorf("Channel already exist")
	}

	err := session.DB(dbName).C("channels").Insert(channel)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Server error, try again")
	}

	return channel, err
}

// GetChannelsForBots returns list of channel for list of bots
func GetChannelsForBots(app *state.AppState, bots []Bot) ([]Channel, error) {
	ids := []bson.ObjectId{}
	for _, bot := range bots {
		ids = append(ids, bot.ID)
	}
	session := app.MgoSession.Clone()
	defer session.Close()

	var channels []Channel

	err := session.DB(dbName).C("channels").Find(bson.M{"botId": bson.M{"$in": ids}}).All(&channels)
	if err != nil {
		log.Println(err)
	}

	return channels, err
}

func GetChannelByID(app *state.AppState, channelID bson.ObjectId) (*Channel, error) {
	session := app.MgoSession.Clone()
	defer session.Close()

	var channel Channel

	err := session.DB(dbName).C("channels").FindId(channelID).One(&channel)
	if err != nil {
		log.Println(err)
	}

	return &channel, err
}

func GetChannelByTelegramID(app *state.AppState, channelID int64) (*Channel, error) {
	session := app.MgoSession.Clone()
	defer session.Close()

	var channel Channel

	err := session.DB(dbName).C("channels").Find(bson.M{"id": channelID}).One(&channel)
	if err != nil {
		log.Println(err)
	}

	return &channel, err
}
