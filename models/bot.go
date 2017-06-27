package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/slonoed/chantra/state"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	tg "gopkg.in/telegram-bot-api.v4"
)

// Bot represents bot in mongo
type Bot struct {
	Token      string        `bson:"token" json:"-"`
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TelegramID int           `json:"telegramId" bson:"telegramId"`
	UserName   string        `json:"username" bson:"username"`
	FirstName  string        `json:"firstName" bson:"firstName"`
	LastName   string        `json:"lastName" bson:"lastName"`

	UserID bson.ObjectId `bson:"userId" json:"-"`
}

// CollectUpdates s
func (bot *Bot) CollectUpdates(app *state.AppState) {
	api, err := tg.NewBotAPI(bot.Token)
	if err != nil {
		log.Println(err)
		return
	}
	u := tg.NewUpdate(0)
	u.Timeout = 1000
	updates, err := api.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		return
	}

	mgoSession := app.MgoSession.Clone()
	defer mgoSession.Close()
	for {
		select {
		case tgUpdate := <-updates:
			update := CreateUpdate(tgUpdate, bot)
			update.Save(mgoSession)
			updateAnswer(app, update, bot)
		case <-app.Done:
			fmt.Println("Done with bot updates " + bot.FirstName)
			return
		}
	}
}

// GetBotByToken returns new or exiting bot, update bot token if needed
func GetBotByToken(app *state.AppState, token string, userID bson.ObjectId) (*Bot, error) {
	session := app.MgoSession.Clone()
	defer session.Close()

	// get from telegram
	botAPI, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	tgUser, err := botAPI.GetMe()
	if err != nil {
		return nil, err
	}

	bot := Bot{
		TelegramID: tgUser.ID,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
		UserName:   tgUser.UserName,
		Token:      token,
		UserID:     userID,
	}
	// find in DB
	coll := session.DB(dbName).C("bots")
	info, err := coll.Upsert(bson.M{"telegramId": bot.TelegramID}, bot)
	if err != nil {
		return nil, err
	}
	if info.UpsertedId == nil {
		err = coll.Find(bson.M{"telegramId": bot.TelegramID}).One(&bot)
	} else {
		err = coll.FindId(info.UpsertedId).One(&bot)
		go bot.CollectUpdates(app)
	}
	if err != nil {
		return nil, err
	}

	return &bot, nil
}

// GetBotsForUser returns list of bots by user rules
func GetBotsForUser(session *mgo.Session, user *User) ([]Bot, error) {
	var bots []Bot
	err := session.DB(dbName).C("bots").Find(bson.M{"userId": user.ID}).All(&bots)
	return bots, err
}

// GetAllBots returns list of all bots
func GetAllBots(session *mgo.Session) ([]Bot, error) {
	var bots []Bot
	err := session.DB(dbName).C("bots").Find(nil).All(&bots)
	return bots, err
}

func GetBotById(app *state.AppState, botID string) (Bot, error) {
	session := app.MgoSession.Clone()
	defer session.Close()
	var bot Bot
	err := session.DB(dbName).C("bots").FindId(bson.ObjectIdHex(botID)).One(&bot)
	return bot, err
}

func updateAnswer(app *state.AppState, update Update, bot *Bot) {
	if update.CallbackQuery == nil {
		return
	}
	var fields []string
	data := update.CallbackQuery.Data
	fields = strings.Split(data, ":")
	if len(fields) != 3 {
		log.Println("Wrong callback fields length", data)
		return
	}
	if !bson.IsObjectIdHex(fields[0]) {
		log.Println("Data post id malformed", fields[0])
		return
	}
	postID := bson.ObjectIdHex(fields[0])
	column, err := strconv.Atoi(fields[2])
	post, err := GetPostById(app, postID)
	if err != nil {
		log.Println("Can't find post", postID, fields[0], err)
		return
	}

	voted, err := update.IsAlreadyVoted(app)
	if err != nil {
		log.Println(err)
		return
	}
	api, err := tg.NewBotAPI(bot.Token)
	if err != nil {
		log.Println(err)
		return
	}
	if voted {
		cb := tg.NewCallback(update.CallbackQuery.ID, "Вы уже проголосовали!")
		cb.CacheTime = 1000 * 60 * 60 * 24 * 365
		api.AnswerCallbackQuery(cb)
		return
	}

	if err = post.AddVote(app, column); err != nil {
		log.Println(err)
	} else {
		cb := tg.NewCallback(update.CallbackQuery.ID, "Ваш голос учтён!")
		go api.AnswerCallbackQuery(cb)
		msg := tg.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			post.BuildReplyMarkup(),
		)
		api.Send(msg)
	}
}
