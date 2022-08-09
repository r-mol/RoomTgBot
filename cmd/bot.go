package main

import (
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Setup() {
	pref := tele.Settings{
		Token:  "5455937729:AAEVDvLDJczTncZ0aOfIA0Xn6dVcFgcMIO0",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	fmt.Println("5455937729:AAEVDvLDJczTncZ0aOfIA0Xn6dVcFgcMIO0")

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
