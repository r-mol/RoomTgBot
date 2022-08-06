package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Setup() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	fmt.Println(os.Getenv("TOKEN"))
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	bot.Handle("/start", func(ctx tele.Context) error {
		return ctx.Send("Start command handler")
	})
	bot.Handle("/help", func(ctx tele.Context) error {
		return ctx.Send("Help command handler")
	})
	bot.Handle(tele.OnText, func(ctx tele.Context) error {
		return ctx.Send(ctx.Message().Text)
	})

	bot.Start()
}
