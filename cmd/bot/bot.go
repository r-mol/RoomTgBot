package bot

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/mongodb"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

var rdb *redis.Client
var mdb *mongo.Client
var mu sync.Mutex

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	rdb.Ping(contex)
}

func init() {
	var err error
	mdb, err = mongodb.NewClient()

	if err != nil {
		panic(err)
	}

	err = mongodb.Ping(mdb)
	if err != nil {
		panic(fmt.Errorf("Ping to MongoDB is unsuccessful: %v", err))
	}
}

func Setup() {
	pref := telegram.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &telegram.LongPoller{Timeout: consts.TimeOutMultiplier * time.Second},
	}

	bot, err := telegram.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	ticker := time.NewTicker(time.Minute)
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				mu.Lock()
				users := map[int64]telegram.User{}

				err = user.GetUsersFromDB(contex, rdb, users)
				if err != nil {
					log.Println(err)
					break
				}

				for key := range users {
					u := users[key]

					err = state.CheckUserOnAvailableNotifications(contex, bot, &u, rdb)
					if err != nil && err != redis.Nil {
						log.Println(err)
						break
					}
				}

				mu.Unlock()
			default:
				continue
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		mu.Lock()
		handling(bot, rdb, mdb)

		mu.Unlock()
	}()

	bot.Start()
	wg.Wait()
}
