package models

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const dbName = "chantra"

// User represent user entity in system
type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email     string        `json:"email" bson:"email"`
	Name      string        `json:"name" bson:"name"`
	Picture   string        `json:"picture" bson:"picture"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
}

// GoogleProfile user profile from google auth
type GoogleProfile struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// UserFromProfile return new or existing user for google profile (by email)
func UserFromProfile(session *mgo.Session, gprofile GoogleProfile) (*User, error) {
	coll := session.DB(dbName).C("users")
	user := User{
		Email:     gprofile.Email,
		Name:      gprofile.Name,
		Picture:   gprofile.Picture,
		CreatedAt: time.Now(),
	}
	info, err := coll.Upsert(bson.M{"email": gprofile.Email}, user)
	if err != nil {
		return nil, err
	}
	if info.UpsertedId == nil {
		err = coll.Find(bson.M{"email": gprofile.Email}).One(&user)
	} else {
		err = coll.FindId(info.UpsertedId).One(&user)
	}
	if err != nil {
		return nil, err
	}

	return &user, err
}

// UserByID get used by ID from DB
func UserByID(session *mgo.Session, id string) (*User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Not bson.ObjectID")
	}
	coll := session.DB(dbName).C("users")
	var user User
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	return &user, err
}
