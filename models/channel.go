package models

import (
	"bitbucket.org/slonoed/chntr/state"
	"gopkg.in/mgo.v2/bson"
	tg "gopkg.in/telegram-bot-api.v4"
)

// Channel represents channel chat in mongo
type Channel struct {
	tg.Chat
	Bots []bson.ObjectId `bson:"bots" json:"bots"`
}

// GetChannelsForBots returns list of channel (uniq) for list of bots
func GetChannelsForBots(app *state.AppState, bots []Bot) ([]Channel, error) {
	// collect bots ids
	ids := []bson.ObjectId{}
	for _, bot := range bots {
		ids = append(ids, bot.ID)
	}
	session := app.MgoSession.Clone()
	defer session.Close()

	// query all updates from db
	var updates []Update
	// TODO move to updates?
	err := session.DB(dbName).C("updates").Find(bson.M{"botId": bson.M{"$in": ids}}).All(&updates)
	if err != nil {
		return nil, err
	}

	// collect uniq channels, if find duplicate, add bot id
	channels := make(map[int64]Channel) // channel id -> channel
	for _, update := range updates {
		post := update.ChannelPost
		if post == nil {
			continue
		}
		chat := post.Chat
		if chat == nil {
			continue
		}

		channel, exist := channels[chat.ID]
		if exist {
			channel.Bots = append(channel.Bots, update.BotID)
		} else {
			channel = Channel{*chat, []bson.ObjectId{update.BotID}}
			channels[chat.ID] = channel
		}
	}

	// convert map to slice
	list := make([]Channel, 0, len(channels))
	for _, channel := range channels {
		list = append(list, channel)
	}

	return list, nil
}
