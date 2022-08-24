package types

import (
	"RoomTgBot/internal/state"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID primitive.ObjectID

type MongoObject interface {
	User | ShoppingEntry | Activity | ExamEntry
}

// -------- Users -----------------------

type User struct {
	MongoID    ID    `json:"_id" bson:"_id,omitempty"`
	TelegramID int64 `json:"telegram_id" bson:"telegram_id"`

	TelegramUsername string `json:"telegram_username" bson:"telegram_username"`
	FirstName        string `json:"first_name" bson:"first_name"`

	NotificationList map[ID]bool `json:"notification_list" bson:"notification_list"`
	ScoreList        map[ID]int  `json:"score_list" bson:"score_list"`

	Order    uint `json:"order" bson:"order"`
	IsAbsent bool `json:"is_absent" bson:"is_absent"`
	IsBot    bool `json:"is_bot" bson:"is_bot"`
}

// -------- Shopping -----------------------

type ShoppingEntry struct {
	MongoID    ID             `json:"_id",bson:"_id",omitempty`
	Photos     state.Messages `json:"shopping_items",bson:"shopping_items"`
	Bill       state.Message  `json:"bill",bson:"bill"`
	TotalPrice float64        `json:"total_price",bson:"total_price"`
	Person     User           `json:"user",bson:"user"`
	Date       time.Time      `json:"date",bson:"date"`
}

// -------- Activities -----------------------

type Activity struct {
	MongoID          ID        `json:"_id",bson:"_id",omitempty`
	Name             string    `json:"name",bson:"name"`
	ScorePerActivity int       `json:"score_per_activity",bson:"score_per_activity"`
	ScoreMultiplier  int       `json:"score_multiplier",bson:"score_multiplier"`
	Scheduled        time.Time `json:"scheduled",bson:"scheduled"`
	RepeatEach       time.Time `json:"repeat_each",bson:"repeat_each"`
	// peolpe circularQueue <person>
}

// -------- Files -----------------------

type ExamMetaData struct {
	Year     uint   `json:"year",bson:"year"`
	Semester uint   `json:"semester",bson:"semester"`
	Course   string `json:"course",bson:"course"`
	Kind     string `json:"kind",bson:"kind"`
}

type ExamEntry struct {
	MongoID  ID             `json:"_id",bson:"_id",omitempty`
	MetaData ExamMetaData   `json:"meta_data",bson:"meta_data"`
	Files    state.Messages `json:"files",bson:"files"`
}
