package bot

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

const timeOutMultiplier = 10

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Ping(contex)
}

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

	handling(bot, rdb)

	bot.Start()
}
