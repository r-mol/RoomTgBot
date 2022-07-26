package types

import (
    "time"

	telegram "gopkg.in/telebot.v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoObject interface {
	MongoId() primitive.ObjectID
	User | ShoppingEntry | Activity | ExamEntry
}

// -------- Users -----------------------

type User struct {
	MongoID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	TelegramID int64              `json:"telegram_id" bson:"telegram_id"`

	TelegramUsername string `json:"telegram_username" bson:"telegram_username"`
	FirstName        string `json:"first_name" bson:"first_name"`

	NotificationList map[primitive.ObjectID]bool `json:"notification_list" bson:"notification_list"`
	ScoreList        map[primitive.ObjectID]int  `json:"score_list" bson:"score_list"`

	Order    uint `json:"order" bson:"order"`
	IsAbsent bool `json:"is_absent" bson:"is_absent"`
	IsBot    bool `json:"is_bot" bson:"is_bot"`
}

func (user User) MongoId() primitive.ObjectID {
	return user.MongoID
}

// -------- Shopping -----------------------

type ShoppingEntry struct {
	MongoID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Photos     []telegram.Photo    `json:"shopping_items" bson:"shopping_items"`
	Bill       telegram.Photo      `json:"bill" bson:"bill"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	Person     User               `json:"user" bson:"user"`
	Date       time.Time          `json:"date" bson:"date"`
}

func (shoppingEntry ShoppingEntry) MongoId() primitive.ObjectID {
	return shoppingEntry.MongoID
}

// -------- Activities -----------------------

type Activity struct {
	MongoID          primitive.ObjectID `json:"_id" bson:"_id ,omitempty"`
	Name             string             `json:"name" bson:"name"`
	ScorePerActivity int                `json:"score_per_activity" bson:"score_per_activity"`
	ScoreMultiplier  int                `json:"score_multiplier" bson:"score_multiplier"`
	Scheduled        time.Time          `json:"scheduled" bson:"scheduled"`
	RepeatEach       time.Time          `json:"repeat_each" bson:"repeat_each"`
}

func (activity Activity) MongoId() primitive.ObjectID {
	return activity.MongoID
}

// -------- Files -----------------------

type Files struct {
	Text   string              `json:"text"`
	Files  []telegram.Document `json:"document"`
	Photos []telegram.Photo    `json:"photo"`
}

type ExamMetaData struct {
	Year     uint   `json:"year" bson:"year"`
	Semester uint   `json:"semester" bson:"semester"`
	Course   string `json:"course" bson:"course"`
	Kind     string `json:"kind" bson:"kind"`
}

type ExamEntry struct {
	MongoID  primitive.ObjectID `json:"_id" bson:"_id ,omitempty"`
	MetaData ExamMetaData       `json:"meta_data" bson:"meta_data"`
	Files    Files              `json:"files" bson:"files"`
}

func (examEntry ExamEntry) MongoId() primitive.ObjectID {
	return examEntry.MongoID
}
