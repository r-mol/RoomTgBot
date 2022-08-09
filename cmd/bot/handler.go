package bot

import (
	"RoomTgBot/internal"
	telegram "gopkg.in/telebot.v3"
)

const (
	commandStart = "/start"
	commandStop  = "/stop"
	commandHelp  = "/help"
	commandWater = "/bring_water"
	commandClean = "/clean_room"
)

func startHandling(bot *telegram.Bot) {
	bot.Handle(commandStart, func(ctx telegram.Context) error {
		newUser := &internal.User{
			ID:        ctx.Sender().ID,
			FirstName: ctx.Sender().FirstName,
			Username:  ctx.Sender().Username,
			IsBot:     ctx.Sender().IsBot}

		if newUser.IsBot {
			defer bot.Stop()
			return ctx.Send("You are fucking bot...")
		}

		// Add to new user to database

		return ctx.Send("Nice to meet you " + newUser.FirstName + "!!!")
	})

	bot.Handle(commandHelp, func(ctx telegram.Context) error {
		return ctx.Send("Help command handler")
	})

	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		return ctx.Send(ctx.Message().Text)
	})
}
