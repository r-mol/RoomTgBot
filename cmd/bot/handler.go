package bot

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"

	"context"
	"log"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

var contex = context.Background()

func handling(bot *telegram.Bot, rdb *redis.Client) {
	menus.InitializeMenus()

	allMenus := menus.GetMenus()

	bot.Handle(commands.CommandStart, func(ctx telegram.Context) error {
		newUser := &user.User{}
		err := user.CreateUser(bot, ctx, newUser)

		if err != nil {
			return err
		}

		log.Println("User is authorized")

		// TODO Add new user to database

		curState := &state.State{
			InitState: commands.CommandStart,
			IsNow:     true,
		}

		states := state.States{}

		states[commands.CommandStart] = curState
		states[state.InitState] = curState

		err = state.SetStatesToRDB(contex, rdb, ctx, &states)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		log.Println(curState)

		return ctx.Send("Nice to meet you "+newUser.FirstName+" !!!", menus.MainMenu)
	})

	bot.Handle(commands.CommandBringWater, func(ctx telegram.Context) error {
		// TODO Find person in database
		//   tgUser := &user.User{}
		//   tgUser = testUser

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		commandFrom := curState.InitState
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commands.CommandStart)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		// TODO Add new data of user to database
		//   testUser = tgUser

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(commands.CommandCleanRoom, func(ctx telegram.Context) error {
		// TODO Find person in database
		//  tgUser := &user.User{}
		//  tgUser = testUser

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		commandFrom := curState.InitState
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commands.CommandStart)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		// TODO Add new data of user to database
		//  testUser = tgUser

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(&menus.BtnRoom, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandRoom)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the room menu...", menus.RoomMenu)
	})

	bot.Handle(&menus.BtnExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandExam)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the exam menu...", menus.ExamMenu)
	})

	bot.Handle(&menus.BtnNews, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandNews)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send text/file messages to create news:", menus.NewsMenu)
	})

	bot.Handle(&menus.BtnSettings, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandSettings)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the settings menu...", menus.SettingsMenu)
	})

	bot.Handle(&menus.BtnDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandNews, commands.CommandDone)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, check your draft of news:", menus.AcceptNewsMenu)
	})

	bot.Handle(&menus.BtnAquaMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandAquaMan)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are aqua-man...", menus.AquaManMenu)
	})

	bot.Handle(&menus.BtnCleanMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandCleanMan)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are clean-man...", menus.CleanManMenu)
	})

	bot.Handle(&menus.BtnShop, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandShop)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the shop menu...", menus.ShopMenu)
	})

	bot.Handle(&menus.BtnBack, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		if curState.InitState == commands.CommandStart {
			return nil
		}

		commandFrom := curState.InitState
		commandTo := curState.PrevState
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commandTo)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Welcome back 🛑", allMenus[commandTo])
	})

	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		return ctx.Send(ctx.Message().Text)
	})
}
