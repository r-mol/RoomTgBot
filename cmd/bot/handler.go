package bot

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"
	"log"

	"context"
	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

var contex = context.Background()

func handling(bot *telegram.Bot) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Ping(contex)

	menus.MainMenu.Reply(
		menus.MainMenu.Row(menus.BtnRoom, menus.BtnNews),
		menus.MainMenu.Row(menus.BtnExam, menus.BtnSettings),
	)

	menus.RoomMenu.Reply(
		menus.RoomMenu.Row(menus.BtnShop, menus.BtnAquaMan, menus.BtnCleanMan),
		menus.MainMenu.Row(menus.BtnBack),
	)

	menus.AquaManMenu.Reply(
		menus.AquaManMenu.Row(menus.BtnBringWater),
		menus.MainMenu.Row(menus.BtnBack),
	)

	bot.Handle(commands.CommandStart, func(ctx telegram.Context) error {

		newUser := &user.User{}
		err := user.CreateUser(bot, ctx, newUser)
		newUser.CurState = commands.CommandStart

		if err != nil {
			return err
		}

		log.Println("User is authorized")

		// Add to new user to database

		st := &state.State{
			InitState: commands.CommandStart,
			IsNow:     true,
		}

		err = state.SetStateToRDB(rdb, contex, st)

		if err != nil {
			return err
		}
		log.Println(st)

		return ctx.Send("Nice to meet you "+newUser.FirstName+" !!!", menus.MainMenu)
	})

	bot.Handle(commands.CommandBringWater, func(ctx telegram.Context) error {
		tgUser := &user.User{}

		// Find person in database

		tgUser.CurState = commands.CommandStart

		// Add new data of user to database

		err := state.CheckOfUserState(rdb, ctx, contex, commands.CommandAquaMan, commands.CommandStart)

		if err != nil {
			return err
		}

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(commands.CommandClean, func(ctx telegram.Context) error {
		tgUser := &user.User{}

		// Find person in database

		tgUser.CurState = commands.CommandStart

		// Add new data of user to database

		err := state.CheckOfUserState(rdb, ctx, contex, commands.CommandCleanMan, commands.CommandStart)

		if err != nil {
			return err
		}
		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(&menus.BtnRoom, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(rdb, ctx, contex, commands.CommandStart, commands.CommandRoom)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// Find person in database

		tgUser.CurState = commands.CommandRoom

		// Add new data of user to database

		return ctx.Send("Now you are in the room", menus.RoomMenu)
	})

	bot.Handle(&menus.BtnAquaMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(rdb, ctx, contex, commands.CommandRoom, commands.CommandAquaMan)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// Find person in database

		tgUser.CurState = commands.CommandAquaMan

		// Add new data of user to database

		return ctx.Send("Now you are aqua-man", menus.AquaManMenu)
	})

	//bot.Handle(&menus.BtnBack, func(ctx telegram.Context) error {
	//
	//	return ctx.Send("Now you are aqua-man", menus.AquaManMenu)
	//})

	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		return ctx.Send(ctx.Message().Text)
	})
}
