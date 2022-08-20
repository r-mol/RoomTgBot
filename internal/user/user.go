package user

import (
	"RoomTgBot/internal/consts"
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

type User struct {
	ID int64 `json:"id"`

	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

func CreateUser(contex context.Context, rdb *redis.Client, bot *telegram.Bot, ctx telegram.Context) error {
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
	return strconv.FormatInt(u.ID, consts.BaseForConvertToInt)
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
