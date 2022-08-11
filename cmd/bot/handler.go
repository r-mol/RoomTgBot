package bot

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"
	"fmt"

	"context"
	"log"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

var contex = context.Background()

func handling(bot *telegram.Bot, rdb *redis.Client) {

	testUser := &user.User{
		ID: 471895149,

		FirstName: "Roman",
		Username:  "roman_molochkov",
		IsBot:     false,
		CurState:  &state.State{},
	}

	menus.InitializeMenus()

	bot.Handle(commands.CommandStart, func(ctx telegram.Context) error {
		newUser := &user.User{}
		err := user.CreateUser(bot, ctx, newUser)

		if err != nil {
			return err
		}

		log.Println("User is authorized")

		// TODO Add to new user to database

		curState := &state.State{
			InitState: commands.CommandStart,
			IsNow:     true,
		}

		newUser.CurState = curState

		// TODO Add new data of user to database
		testUser = newUser

		err = state.SetStateToRDB(contex, rdb, curState)

		if err != nil {
			return err
		}
		log.Println(curState)

		return ctx.Send("Nice to meet you "+newUser.FirstName+" !!!", menus.MainMenu)
	})

	bot.Handle(commands.CommandBringWater, func(ctx telegram.Context) error {
		curState, err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandAquaMan, commands.CommandStart)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// TODO Find person in database
		tgUser = testUser

		tgUser.CurState = curState

		// TODO Add new data of user to database
		testUser = tgUser

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(commands.CommandClean, func(ctx telegram.Context) error {
		curState, err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandCleanMan, commands.CommandStart)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// TODO Find person in database
		tgUser = testUser

		tgUser.CurState = curState

		// TODO Add new data of user to database
		testUser = tgUser

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(&menus.BtnRoom, func(ctx telegram.Context) error {
		curState, err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandRoom)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// TODO Find person in database
		tgUser = testUser

		tgUser.CurState = curState

		// TODO Add new data of user to database
		testUser = tgUser

		return ctx.Send("Now you are in the room", menus.RoomMenu)
	})

	bot.Handle(&menus.BtnAquaMan, func(ctx telegram.Context) error {
		curState, err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandAquaMan)

		if err != nil {
			return err
		}

		tgUser := &user.User{}

		// TODO Find person in database
		tgUser = testUser

		tgUser.CurState = curState

		// TODO Add new data of user to database
		testUser = tgUser

		return ctx.Send("Now you are aqua-man", menus.AquaManMenu)
	})

	bot.Handle(&menus.BtnBack, func(ctx telegram.Context) error {
		fmt.Println(testUser)
		fmt.Println(testUser.CurState)
		return ctx.Send("salam")
	})

	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		return ctx.Send(ctx.Message().Text)
	})
}
