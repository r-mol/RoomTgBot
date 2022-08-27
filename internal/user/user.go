package user

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/mongodb"
	"RoomTgBot/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	telegram "gopkg.in/telebot.v3"
)

func CreateUser(contex context.Context, rdb *redis.Client, bot *telegram.Bot, ctx telegram.Context) error {
	idString := strconv.FormatInt(ctx.Sender().ID, consts.BaseForConvertToInt)

	_, err := rdb.Get(contex, idString).Result()
	if err == redis.Nil {
		// TODO add new user to mongoDB and don't forget to normalizeOrder
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

// Normalize order of people so the smallest one is 0 and others ared 1, 2...
func normalizeOrder(users []types.User) []types.User {
	sort.Slice(users, func(p, q int) bool {
		return users[p].Order < users[q].Order
	})

	for index := range users {
		users[index].Order = uint(index)
	}

	return users
}

// Users list and map should be normalized using normalizeOrder
func NextInOrder(prevID int64, usersMap map[int64]types.User, users []types.User) (int64, error) {
	prevOrder := usersMap[prevID].Order

	if len(users) == 0 {
		return 0, fmt.Errorf("no next user, because list of users is empty")
	}

	index := (int(prevOrder) + 1) % len(users)

	return users[index].TelegramID, nil
}

func MongoGetMap(ctx context.Context, client *mongo.Client) (map[int64]types.User, error) {
	mongoUsers, err := mongodb.GetAll[types.User](ctx, client, consts.MongoUsersCollection)
	if err != nil {
		return map[int64]types.User{}, fmt.Errorf("unable to get Users from mongodb: %v", err)
	}

	users := map[int64]types.User{}
	for _, elem := range mongoUsers {
		users[elem.TelegramID] = elem
	}

	return users, nil
}
func MongoGet(ctx context.Context, client *mongo.Client) ([]types.User, error) {
	mongoUsers, err := mongodb.GetAll[types.User](ctx, client, "Users")
	if err != nil {
		return []types.User{}, fmt.Errorf("unable to get Users from mongodb: %v", err)
	}

	return mongoUsers, nil
}
func MongoAdd(ctx context.Context, client *mongo.Client, user *types.User) error {
	_, err := mongodb.AddOne(ctx, client, consts.MongoUsersCollection, user)

	if err != nil {
		return fmt.Errorf("unable to add user to mongodb: %v", err)
	}

	return nil
}

// func (u *User) Recipient() string {
// 	return strconv.FormatInt(u.TelegramID, consts.BaseForConvertToInt)
// }

func GetUsersFromDB(contex context.Context, rdb *redis.Client, users map[int64]telegram.User) error {
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

	err := GetUsersFromDB(contex, rdb, users)
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
