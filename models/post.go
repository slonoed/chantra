package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/slonoed/chantra/state"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	tg "gopkg.in/telegram-bot-api.v4"
)

type ParseMode string

type Answer struct {
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// Post represents post in database
type Post struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ChannelID        bson.ObjectId `json:"channel_id" bson:"channelid"`
	IsSent           bool          `json:"isSent" bson:"isSent"`
	SentDate         Timestamp     `json:"sentDate" bson:"sentDate"`
	Text             string        `json:"text"`
	WithNotification bool          `json:"withТotification"`
	WithPreview      bool          `json:"withPreview"`
	Mode             ParseMode     `json:"parseMode"`
	Answers          []Answer      `json:"answers"`
}

func GetPostById(app *state.AppState, id bson.ObjectId) (*Post, error) {
	session := app.CloneSession()
	defer session.Close()
	var post Post
	if err := session.DB(dbName).C("posts").FindId(id).One(&post); err != nil {
		return nil, err
	}
	return &post, nil
}

// Save post to database
func (post *Post) Save(app *state.AppState) error {
	session := app.MgoSession.Clone()
	defer session.Close()
	return session.DB(dbName).C("posts").Insert(post)
}

// Send post to telegram
func (post *Post) Send(app *state.AppState) error {
	session := app.MgoSession.Clone()
	defer session.Close()

	channel, err := GetChannelByID(app, post.ChannelID)
	if err != nil {
		return err
	}

	bot, err := GetBotById(app, channel.BotID.Hex())
	if err != nil {
		return err
	}
	botAPI, err := tg.NewBotAPI(bot.Token)
	if err != nil {
		return err
	}

	coll := session.DB(dbName).C("posts")
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"isSent": true}},
		ReturnNew: false,
	}
	var prevPostState Post
	_, err = coll.FindId(post.ID).Apply(change, &prevPostState)
	if err != nil {
		return err
	}

	// return if post already processed by something
	if prevPostState.IsSent {
		return nil
	}

	msg := tg.NewMessage(channel.Chat.ID, post.Text)
	msg.ReplyMarkup = post.BuildReplyMarkup()
	msg.DisableNotification = !post.WithNotification
	msg.DisableWebPagePreview = !post.WithPreview
	if post.Mode == tg.ModeMarkdown {
		msg.ParseMode = tg.ModeMarkdown
	} else {
		msg.ParseMode = ""
	}

	_, err = botAPI.Send(msg)
	if err != nil {
		if err := coll.UpdateId(post.ID, bson.M{"$set": bson.M{"isSent": false}}); err != nil {
			log.Println(err)
		}
		return err
	}

	return nil
}

func (post *Post) BuildReplyMarkup() tg.InlineKeyboardMarkup {
	row := []tg.InlineKeyboardButton{}
	for idx, answer := range post.Answers {
		data := post.ID.Hex() + ":0:" + strconv.Itoa(idx)
		text := fmt.Sprintf("%s (%d)", answer.Text, answer.Votes)
		btn := tg.NewInlineKeyboardButtonData(text, data)
		row = append(row, btn)
	}
	return tg.NewInlineKeyboardMarkup(row)
}

// GetReadyToSendPosts return posts with sentDate before now
func GetReadyToSendPosts(app *state.AppState) ([]Post, error) {
	session := app.MgoSession.Clone()
	defer session.Close()
	var posts []Post
	query := bson.M{"isSent": false, "sentDate": bson.M{"$lt": Now()}}
	err := session.DB(dbName).C("posts").Find(query).All(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (post *Post) AddVote(app *state.AppState, index int) error {
	s := app.CloneSession()
	defer s.Close()
	// check already voted
	coll := s.DB(dbName).C("posts")
	path := fmt.Sprintf("answers.%d.votes", index)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{path: 1}},
		ReturnNew: true,
	}
	if _, err := coll.FindId(post.ID).Apply(change, post); err != nil {
		return err
	}

	return nil
}

// GetPostsForChannels returns list of posts for list of channels
func GetPostsForChannels(app *state.AppState, channels []Channel) ([]Post, error) {
	// collect bots ids
	var ids []bson.ObjectId
	for _, channel := range channels {
		ids = append(ids, channel.ID)
	}
	session := app.CloneSession()
	defer session.Close()

	var posts []Post
	query := bson.M{"channelid": bson.M{"$in": ids}}
	if err := session.DB(dbName).C("posts").Find(query).All(&posts); err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}
