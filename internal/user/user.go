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

// ---------------------------Score-------------------------------------

func IncreaseScore(tgID int64, usersMap map[int64]types.User, activityName string,
	activityMap map[string]types.Activity) (map[int64]types.User, error) {
	return changeScore(tgID, usersMap, activityName, activityMap, 1)
}

func DecreaseScore(tgID int64, usersMap map[int64]types.User, activityName string,
	activityMap map[string]types.Activity) (map[int64]types.User, error) {
	return changeScore(tgID, usersMap, activityName, activityMap, -1)
}
func changeScore(tgID int64, usersMap map[int64]types.User, activityName string,
	activityMap map[string]types.Activity, sign int) (map[int64]types.User, error) {
	activity, ok := activityMap[activityName]
	if !ok {
		return usersMap, fmt.Errorf("such activity does not exist: %s", activityName)
	}

	user, ok := usersMap[tgID]
	if !ok {
		return usersMap, fmt.Errorf("user with id %v does not exist", tgID)
	}

	user.ScoreList[activity.MongoID] += activity.ScoreMultiplier * activity.ScorePerActivity * sign

	return usersMap, nil
}

// ---------------------------Order-------------------------------------

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

func ChangeOrder(ctx context.Context, client *mongo.Client, users []types.User, indexIDmap map[int64]uint) error {
	updatedUsers := []types.User{}

	for _, user := range users {
		user.Order = indexIDmap[user.TelegramID]
		updatedUsers = append(updatedUsers, user)
	}

	updatedUsers = normalizeOrder(updatedUsers)

	err := mongodb.UpdateAll(ctx, client, consts.MongoUsersCollection, updatedUsers)
	if err != nil {
		return fmt.Errorf("unable to change order of users: %v", err)
	}

	// TODO update order of users in Redis
	return nil
}

// Next order value of a person
func nextOrderValue(users []types.User) int {
	return len(users)
}

// Users list and map should be normalized using normalizeOrder
func NextInOrder(prevID int64, usersMap map[int64]types.User, users []types.User) (int64, error) {
	prevOrder := usersMap[prevID].Order

	if len(users) == 0 {
		return 0, fmt.Errorf("no next user: list of users is empty")
	}

	index := (int(prevOrder) + 1) % len(users)

	return users[index].TelegramID, nil
}

// ---------------------------Databases-------------------------------------

// Get map of [telegramID]user from Mongodb
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

// Get map of [activityName]activity from Mongodb
func MongoActivitiesMap(ctx context.Context, client *mongo.Client) (map[string]types.Activity, error) {
	mongoActivities, err := mongodb.GetAll[types.Activity](ctx, client, consts.MongoActivitiesCollection)
	if err != nil {
		return map[string]types.Activity{}, fmt.Errorf("unable to get Users from mongodb: %v", err)
	}

	activities := map[string]types.Activity{}
	for _, elem := range mongoActivities {
		activities[elem.Name] = elem
	}

	return activities, nil
}

// Get user from Mongodb
func MongoGet(ctx context.Context, client *mongo.Client) ([]types.User, error) {
	mongoUsers, err := mongodb.GetAll[types.User](ctx, client, "Users")
	if err != nil {
		return []types.User{}, fmt.Errorf("unable to get Users from mongodb: %v", err)
	}

	return mongoUsers, nil
}

// Add user to Mongodb
func MongoAdd(ctx context.Context, client *mongo.Client, user *types.User) error {
	_, err := mongodb.AddOne(ctx, client, consts.MongoUsersCollection, user)

	if err != nil {
		return fmt.Errorf("unable to add user to mongodb: %v", err)
	}

	return nil
}

// Update user in Mongodb
func MongoUpdate(ctx context.Context, client *mongo.Client, user types.User) error {
	return mongodb.UpdateOne(ctx, client, consts.MongoUsersCollection, user)
}

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
