package bot

import (
	"fmt"
	"log"
	"os"
	"time"

	telegram "gopkg.in/telebot.v3"
)

func Setup() {
	pref := telegram.Settings{
		Token: os.Getenv("TG_TOKEN"),
		Poller: &telegram.LongPoller{Timeout: 10 * time.Second},
	}

	fmt.Println(os.Getenv("TG_TOKEN"))

	bot, err := telegram.NewBot(pref)

	if err != nil {
		log.Fatal(err)
		return
	}

	startHandling(bot)

	bot.Start()
}
