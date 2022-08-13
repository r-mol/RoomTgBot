package bot

import (
	"log"
	"time"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

const timeOutMultiplier = 10

func Setup() {
	pref := telegram.Settings{
		Token:  "5455937729:AAEVDvLDJczTncZ0aOfIA0Xn6dVcFgcMIO0",
		Poller: &telegram.LongPoller{Timeout: timeOutMultiplier * time.Second},
	}

	bot, err := telegram.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Ping(contex)

	handling(bot, rdb)

	bot.Start()
}
