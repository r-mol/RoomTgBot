package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	telegram "gopkg.in/telebot.v3"
	"time"
)

type ID int64

// -------- People -----------------------

// -------- Shopping -----------------------

type ShoppingEntry struct {
	MongoID    ID              `json:"_id",bson:"_id",omitempty`
	Photos     []telegram.File `json:"shopping_items",bson:"shopping_items"`
	Bill       telegram.File   `json:"bill",bson:"bill"`
	TotalPrice float64         `json:"total_price",bson:"total_price"`
	Person     User            `json:"user",bson:"user"`
	Date       time.Time       `json:"date",bson:"date"`
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
	MongoID  ID              `json:"_id",bson:"_id",omitempty`
	MetaData ExamMetaData    `json:"meta_data",bson:"meta_data"`
	Files    []telegram.File `json:"files",bson:"files"`
}
