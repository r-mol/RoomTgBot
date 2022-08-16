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

	handlingRoomMenu(bot, rdb)
	handlingNewsMenu(bot, rdb)
	handlingExamMenu(bot, rdb)
	handlingSettingsMenu(bot, rdb)

	bot.Handle(commands.CommandStart, func(ctx telegram.Context) error {
		newUser := &user.User{}
		err := user.CreateUser(bot, ctx, newUser)

		if err != nil {
			return err
		}

		log.Println("User is authorized")

		// TODO Add new user to database

		curState := &state.State{
			StateName: commands.CommandStart,
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

	bot.Handle(&menus.BtnPrevious, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableListStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("Please restart bot ✨")
		}

		message := curState.GetPrevMessageOfList()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, message)
		if err != nil {
			return err
		}

		return ctx.Send("You can list in the items:", menus.ListMenu)
	})

	bot.Handle(&menus.BtnNext, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableListStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("Please restart bot ✨")
		}

		message := curState.GetNextMessageOfList()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, message)
		if err != nil {
			return err
		}

		return ctx.Send("You can list in the items:", menus.ListMenu)
	})

	bot.Handle(&menus.BtnExit, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return ctx.Send("You exit from list", allMenus[curState.StateName])
	})

	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableChattingStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("You can not write here or you send unavailable command...")
		}

		curState.Text += " " + ctx.Message().Text

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(telegram.OnDocument, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableChattingStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("You can not write here or you send unavailable command...")
		}

		curState.Files = append(curState.Files, ctx.Message().Document)

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(telegram.OnPhoto, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableChattingStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("You can not write here or you send unavailable command...")
		}

		curState.Photos = append(curState.Photos, ctx.Message().Photo)

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(&menus.BtnBack, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		if curState.StateName == commands.CommandStart {
			return nil
		}

		commandFrom := curState.StateName
		commandTo := curState.PrevState
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commandTo)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return ctx.Send("Welcome back 🛑", allMenus[commandTo])
	})
}

func handlingRoomMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnRoom, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandRoom)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the room menu...", menus.RoomMenu)
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

	bot.Handle(&menus.BtnUploadPurchase, func(ctx telegram.Context) error {

		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandShop, commands.CommandUploadPurchase)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send:\n1. Name of purchase\n2. Image of purchase\n3. Cost of purchase", menus.ShopUploadMenu)
	})

	bot.Handle(&menus.BtnCheckShopping, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandShop, commands.CommandCheck)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableListStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("Please restart bot ✨")
		}

		// TODO Get list of purchases from database
		setOfMessages := []*state.Message{{Text: "Test1"}, {Text: "Test2"}, {Text: "Test3"}}
		curState.ListMessage = setOfMessages

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		err = ctx.Send("Here you can check purchases from new to old one:", menus.ShopCheckMenu)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, setOfMessages[0])
		if err != nil {
			return err
		}

		return ctx.Send("You can list in the items:", menus.ListMenu)
	})

	bot.Handle(&menus.BtnPurchaseDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandUploadPurchase, commands.CommandPurchaseDone)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Please, check your draft of purchase:", menus.PostPurchaseMenu)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, nil)
		if err != nil {
			return err
		}

		return err
	})

	bot.Handle(&menus.BtnPostPurchase, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		// TODO add purchase to database and make notification

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		err = state.CheckOfUserState(contex, rdb, ctx, commands.CommandPurchaseDone, commands.CommandStart)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Your purchase has been sand to database 📨", menus.MainMenu)
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

		commandFrom := curState.StateName
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

		commandFrom := curState.StateName
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
}

func handlingNewsMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnNews, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandNews)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send text/file messages to create news:", menus.NewsMenu)
	})

	bot.Handle(&menus.BtnNewsDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandNews, commands.CommandNewsDone)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Please, check your draft of news:", menus.PostNewsMenu)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, nil)
		if err != nil {
			return err
		}

		return err
	})

	bot.Handle(&menus.BtnPostNews, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		// TODO add news to database and make notification

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		err = state.CheckOfUserState(contex, rdb, ctx, commands.CommandNewsDone, commands.CommandStart)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Your news has been posted 📨", menus.MainMenu)
	})

	bot.Handle(&menus.BtnDeleteDraft, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx)
		if err != nil {
			return err
		}

		return ctx.Send("All you messages have been removed")
	})
}

func handlingExamMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandExam)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the exam menu...", menus.ExamMenu)
	})

	bot.Handle(&menus.BtnUploadExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandExam, commands.CommandUploadExam)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send files:", menus.ExamUploadMenu)
	})

	bot.Handle(&menus.BtnExamDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandUploadExam, commands.CommandExamDone)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}
		// TODO Create moving in list of subjects
		err = ctx.Send("Please, check your files of exam and choose subject from list:")
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, nil)
		if err != nil {
			return err
		}

		return err
	})
}
func handlingSettingsMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnSettings, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandStart, commands.CommandSettings)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the settings menu...", menus.SettingsMenu)
	})
}
