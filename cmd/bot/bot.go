package bot

import (
	"fmt"
	"log"
	"time"

	telegram "gopkg.in/telebot.v3"
)

const timeOutMultiplier = 10

func Setup() {
	pref := telegram.Settings{
		Token:  "5455937729:AAEVDvLDJczTncZ0aOfIA0Xn6dVcFgcMIO0",
		Poller: &telegram.LongPoller{Timeout: timeOutMultiplier * time.Second},
	}

	fmt.Println("5455937729:AAEVDvLDJczTncZ0aOfIA0Xn6dVcFgcMIO0")

	bot, err := telegram.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	handling(bot)

	bot.Start()
}
