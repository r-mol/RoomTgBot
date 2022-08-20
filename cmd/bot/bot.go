package bot

import (
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"

	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

const timeOutMultiplier = 10

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	rdb.Ping(contex)
}

var mu sync.Mutex

func Setup() {
	pref := telegram.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &telegram.LongPoller{Timeout: timeOutMultiplier * time.Second},
	}

	bot, err := telegram.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	ticker := time.NewTicker(time.Minute)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				users := map[int64]telegram.User{}
				err = user.GetUserUsersFromDB(contex, rdb, users)
				if err != nil {
					log.Println(err)
				}
				for _, user := range users {
					err = state.CheckUserOnAvaliableNotifications(bot, &user, contex, rdb)
					if err != nil && err != redis.Nil {
						log.Println(err)
					}
				}
				mu.Unlock()
			default:
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		handling(bot, rdb)

		mu.Unlock()
	}()

	bot.Start()
	wg.Wait()
}
