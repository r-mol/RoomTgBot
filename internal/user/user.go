package user

import (
	"RoomTgBot/internal/consts"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	telegram "gopkg.in/telebot.v3"
)

type ID primitive.ObjectID

type User struct {
	MongoID    ID     `json:"_id",bson:"_id",omitempty`
	TelegramID string `json:"telegram_id",bson:"telegram_id"`

	TelegramUsername string `json:"telegram_username",bson:"telegram_username"`
	FirstName        string `json:"first_name",bson:"first_name"`

	NotificationList map[ID]bool `json:"notification_list",bson:"notification_list"`
	ScoreList        map[ID]int  `json:"score_list",bson:"score_list"`

	Order    uint `json:"order",bson:"order"`
	IsAbsent bool `json:"is_absent",bson:"is_absent"`
	IsBot    bool `json:"is_bot",bson:"is_bot"`
}

// type User struct {
// 	ID int64 `json:"id"`
//
// 	FirstName string `json:"first_name"`
// 	Username  string `json:"username"`
// 	IsBot     bool   `json:"is_bot"`
// }

// Create new User

func CreateUser()  {
	// rdbAddUser()
	// mongodbAddUser()
}

func mongodbAddUser(ctx context.Context, db *mongo.Client, user User) (*mongo.InsertOneResult, error) {
	users := db.Database(consts.MongoDBName).Collection("users")
    insertResult, err := users.InsertOne(ctx, user)
    if err != nil{
        return insertResult, fmt.Errorf("Unable to add new user to MongoDB: %v", err)
    }
    return insertResult, nil
}

func rdbAddUser(contex context.Context, rdb *redis.Client, bot *telegram.Bot, ctx telegram.Context) error {
	idString := strconv.FormatInt(ctx.Sender().ID, consts.BaseForConvertToInt)

	_, err := rdb.Get(contex, idString).Result()
	if err == redis.Nil {
		// TODO add new user to database
		err = SetUserToRDB(contex, rdb, ctx)
		if err != nil {
			return err
		}

		if ctx.Sender().IsBot {
			defer bot.Stop()
			return ctx.Send("You are fucking bot...")
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (u *User) Recipient() string {
	return strconv.FormatInt(u.TelegramID, consts.BaseForConvertToInt)
}

func GetUserUsersFromDB()(){
    // if not in redis
    // add to redis from mongo
    
    // return users
       
}

func GetUserUsersFromDB(contex context.Context, rdb *redis.Client, users map[int64]telegram.User) error {
	stateString, err := rdb.Get(contex, "0").Result()

	switch err {
	case nil:
		err = json.Unmarshal([]byte(stateString), &users)
		if err != nil {
			return err
		}
	case redis.Nil:
	default:
		return err
	}

	return nil
}

func SetUserToRDB(contex context.Context, rdb *redis.Client, ctx telegram.Context) error {
	users := map[int64]telegram.User{}

	err := GetUserUsersFromDB(contex, rdb, users)
	if err != nil {
		return err
	}

	users[ctx.Sender().ID] = *ctx.Sender()
	stateBytes, err := json.Marshal(users)

	if err != nil {
		return err
	}

	err = rdb.Set(contex, "0", stateBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
