package types

import (
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ID string

// -------- People -----------------------
type Person struct {
	MongoID       ID     `bson:"_id,omitempty"`
	TelegramAlias string `bson:"telegramAlias"`
	Nickname      string `bson:"nickname"`

	NotificationList map[ID]bool `bson:"notificationList"`
	ScoreList        map[ID]int  `bson:"scoreList"`

	Order    uint `bson:"order"`
	IsAbsent bool `bson:"isAbsent"`
}

// -------- Shopping -----------------------

type ShoppingItem struct {
	Name  string `bson:"name" `
	Photo File   `bson:"photo"`
}
type ShoppingEntry struct {
	MongoID       ID             `bson:"_id,omitempty"`
	ShoppingItems []ShoppingItem `bson:"shoppingItems"`
	Bill          File           `bson:"bill"       `
	TotalPrice    float64        `bson:"totalPrice"  `
	Person        Person         `bson:"person"       `
	Date          time.Time      `bson:"date"      `
}

// -------- Activities -----------------------

type Activity struct {
	MongoID          ID        `bson:"_id, omitempty"`
	Name             string    `bson:"name"          `
	ScorePerActivity int       `bson:"scorePerActivity"`
	ScoreMultiplier  int       `bson:"scoreMultiplier"`
	Scheduled        time.Time `bson:"scheduled"      `
	RepeatEach       time.Time `bson:"repeatEach"     `
	// peolpe circularQueue <person>
}

// -------- Notifications -----------------------

//
// type Notifiable interface {
// 	Notify() error
// }
//
// func (textInformaiton TextInformation) Notify() error {
// 	// TODO: implement notify function
// 	return nil
// }
//
// func (activity Activity) Notify() error {
// 	// TODO: implement notify function
// 	return nil
// }

// -------- Files -----------------------

type TextInformation struct {
	Header string `bson:"header"`
	Body   string `bson:"body"`
}

type File interface {
	// TODO: implement File
}

type FileInformation struct {
	Year     uint            `bson:"year"    `
	Semester uint            `bson:"semester"`
	Course   string          `bson:"course"`
	Kind     string          `bson:"kind"  `
	Info     TextInformation `bson:"info" `
}

type FileEntry struct {
	MongoID  ID              `bson:"_id, omitempty"`
	MetaData FileInformation `bson:"metaData"`
	Files    []File          `bson:"files" `
}
