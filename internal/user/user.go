package user

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/mongodb"
	"RoomTgBot/internal/types"
	"context"
	// "encoding/json"
	"fmt"
	"sort"
	// "strconv"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	telegram "gopkg.in/telebot.v3"
)

type Storage struct {
	RedisClient  *redis.Client
	RedisContext context.Context

	MongoClient  *mongo.Client
	MongoContext context.Context

	Bot        *telegram.Bot
	BotContext telegram.Context

	// usersList []types.User
	// UsersMap  map[int64]*types.User
	// UsersList []*types.User
}

func InitializeStorage(
	redisContext context.Context,
	redisClient *redis.Client,
	mongoContext context.Context,
	mongoClient *mongo.Client,
	botContext telegram.Context,
	bot *telegram.Bot) (*Storage, error) {

	storage := &Storage{
		RedisClient:  redisClient,
		RedisContext: redisContext,
		MongoClient:  mongoClient,
		MongoContext: mongoContext,
		Bot:          bot,
		BotContext:   botContext,
	}

	err := storage.MongoLoad()
	if err != nil {
		return storage, fmt.Errorf("unable to initialize user storage: %v", err)
	}

	return storage, nil
}

func (storage *Storage) CreateUser() error {
	newUser := storage.BotGetUser()
	_, err := storage.redisGet(newUser.TelegramID)

	if err == nil {
		storage.BotContext.Send("User already exists")
		return fmt.Errorf("user already exists")
	}

    if newUser.IsBot {
		storage.BotContext.Send("You are bot")
		return fmt.Errorf("user is bot")
    }

	storage.redisAdd(&newUser)

	return nil
}

// ---------------------------Score-------------------------------------

func (storage *Storage) IncreaseScore(tgID int64, activity *types.Activity) error {
	return storage.changeScore(tgID, activity, 1)
}

func (storage *Storage) DecreaseScore(tgID int64, activity *types.Activity) error {
	return storage.changeScore(tgID, activity, -1)
}

func (storage *Storage) changeScore(tgID int64, activity *types.Activity, sign int) error {
	user, ok := storage.redisGet(tgID)
	if !ok {
		return fmt.Errorf("user with id %v does not exist", tgID)
	}

	user.ScoreList[activity.MongoID] += activity.ScoreMultiplier * activity.ScorePerActivity * sign

	return nil
}

// ---------------------------Order-------------------------------------

// Normalize order of people so the smallest one is 0 and others ared 1, 2...
func (storage *Storage) normalizeOrder() error {
    users, err := storage.redisGetList()
    if err != nil{
        return fmt.Errorf("unable to normalizeOrder: %v", err)
    }

	sort.Slice(users, func(p, q int) bool {
		return users[p].Order < users[q].Order
	})

	for index := range users {
        storage.redisUpdate(users[index].ID, "order", uint(index))
	}
    return nil
}

func (storage *Storage) ChangeOrder(indexIDmap map[int64]uint) error {
	oldUsersList := storage.UsersList

	for _, user := range oldUsersList {
		storage.UsersMap[user.TelegramID].Order = indexIDmap[user.TelegramID]
	}

	storage.normalizeOrder()

	err := mongodb.UpdateAll(storage.MongoContext, storage.MongoClient, consts.MongoUsersCollection, storage.usersList)
	if err != nil {
		return fmt.Errorf("unable to change order of storage.UsersList: %v", err)
	}

	// TODO update order of storage.UsersList in Redis
	return nil
}

// Users list and map should be normalized using normalizeOrder
func (storage *Storage) NextInOrder(prevID int64) (int64, error) {
	prevOrder := storage.UsersMap[prevID].Order
	numberOfUsers := len(storage.UsersList)

	if numberOfUsers == 0 {
		return 0, fmt.Errorf("no next user: list of users is empty")
	}

	index := (int(prevOrder) + 1) % numberOfUsers

	return storage.UsersList[index].TelegramID, nil
}

// ---------------------------Databases-------------------------------------

func (storage *Storage) MongoLoad() error {
	users, err := mongodb.GetAll[types.User](storage.MongoContext, storage.MongoClient, consts.MongoUsersCollection)
	if err != nil {
		return fmt.Errorf("unable to get Users from mongodb: %v", err)
	}

	usersMap := map[int64]*types.User{}
	usersList := []*types.User{}

	for _, elem := range users {
		usersMap[elem.TelegramID] = &elem
		usersList = append(usersList, &elem)
	}

	storage.UsersList = usersList
	storage.UsersMap = usersMap
	storage.usersList = users

	return nil
}

// TODO: move to activities package
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

// func (storage *Storage)Add(user *types.User) error{
// // add to mongo and redis
//
// }

func (storage *Storage) storageAdd(user types.User) {
	storage.usersList = append(storage.usersList, user)
	userPointer := &storage.usersList[len(storage.usersList)-1]
	storage.UsersList = append(storage.UsersList, userPointer)
	storage.UsersMap[user.TelegramID] = userPointer
}

func (storage *Storage) mongoAdd(user *types.User) error {
	_, err := mongodb.AddOne(storage.MongoContext, storage.MongoClient, consts.MongoUsersCollection, user)
	if err != nil {
		return fmt.Errorf("unable to add user to mongodb: %v", err)
	}

	return nil
}

func (storage *Storage) UpdateFromRedis() error {
	// TODO: implement function for updating storage.userList from Redis
	return nil
}

func (storage *Storage) BotGetUser() types.User {
	tgUser := storage.BotContext.Sender()
	return types.User{
		TelegramID:       tgUser.ID,
		TelegramUsername: tgUser.Username,
		FirstName:        tgUser.FirstName,
		IsBot:            tgUser.IsBot,
	}
}


func (storage *Storage) redisAdd(user *types.User) error {
	// TODO: implement only adding user to redis
	return nil
}

func (storage *Storage) redisGet(telegramID int64) (types.User, error) {
	// TODO: implement getting one user by telegram id
}

func (storage *Storage) redisGetList() ([]types.User, error) {
	// TODO: implement getting list of all users
}

func (storage *Storage) redisUpdate(telegramID int64, propertyName string, ) (error) {
	// TODO: implement getting list of all users
}
