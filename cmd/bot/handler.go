package bot

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/exam"
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

	handlingStart(bot, rdb)
	handlingList(bot, rdb, allMenus)
	handlingRoomMenu(bot, rdb)
	handlingNewsMenu(bot, rdb)
	handlingExamMenu(bot, rdb)
	handlingSettingsMenu(bot, rdb)
	handlingTriggersOnMessages(bot, rdb)
}

func handlingStart(bot *telegram.Bot, rdb *redis.Client) {
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
}

func handlingList(bot *telegram.Bot, rdb *redis.Client, allMenus map[string]*telegram.ReplyMarkup) {
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

	handlingShopMenu(bot, rdb)
	handlingAquaMan(bot, rdb)
	handlingCleanMan(bot, rdb)
}

func handlingShopMenu(bot *telegram.Bot, rdb *redis.Client) {
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
}

func handlingAquaMan(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnAquaMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandAquaMan)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are aqua-man...", menus.AquaManMenu)
	})

	bot.Handle(commands.CommandBringWater, func(ctx telegram.Context) error {
		// TODO Find person in database
		//   tgUser := &user.User{}
		//   tgUser = testUser

		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		// TODO Add new data of user to database
		//   testUser = tgUser

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})
}

func handlingCleanMan(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnCleanMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandRoom, commands.CommandCleanMan)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are clean-man...", menus.CleanManMenu)
	})

	bot.Handle(commands.CommandCleanRoom, func(ctx telegram.Context) error {
		// TODO Find person in database
		//  tgUser := &user.User{}
		//  tgUser = testUser

		err := state.ReturnToStartState(contex, rdb, ctx)
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
	handlingSubjects(bot, rdb)

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

	bot.Handle(&menus.BtnGetExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, commands.CommandExam, commands.CommandGetExam)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, choose subject from list:", menus.SubjectMenu)
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

		err = ctx.Send("Please, check your files of exam:")
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessages(ctx, nil)
		if err != nil {
			return err
		}

		err = ctx.Send("And choose subject from list:", menus.SubjectMenu)
		if err != nil {
			return err
		}

		return ctx.Delete()
	})
}

func handlingSubjects(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.Subject1, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject1.Text)
	})

	bot.Handle(&menus.Subject2, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject2.Text)
	})

	bot.Handle(&menus.Subject3, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject3.Text)
	})

	bot.Handle(&menus.Subject4, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject4.Text)
	})

	bot.Handle(&menus.Subject5, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject5.Text)
	})

	bot.Handle(&menus.Subject6, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject6.Text)
	})

	bot.Handle(&menus.Subject7, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject7.Text)
	})

	bot.Handle(&menus.Subject8, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject8.Text)
	})

	bot.Handle(&menus.Subject9, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject9.Text)
	})

	bot.Handle(&menus.Subject10, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject10.Text)
	})

	bot.Handle(&menus.Subject11, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject11.Text)
	})

	bot.Handle(&menus.Subject12, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject12.Text)
	})

	bot.Handle(&menus.Subject13, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject13.Text)
	})

	bot.Handle(&menus.Subject14, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject14.Text)
	})

	bot.Handle(&menus.Subject15, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject15.Text)
	})

	bot.Handle(&menus.Subject16, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject16.Text)
	})

	bot.Handle(&menus.Subject17, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject17.Text)
	})

	bot.Handle(&menus.Subject18, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject18.Text)
	})

	bot.Handle(&menus.Subject19, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject19.Text)
	})

	bot.Handle(&menus.Subject20, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject20.Text)
	})

	bot.Handle(&menus.Subject21, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject21.Text)
	})

	bot.Handle(&menus.Subject22, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject22.Text)
	})

	bot.Handle(&menus.Subject23, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject23.Text)
	})

	bot.Handle(&menus.Subject24, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject24.Text)
	})

	bot.Handle(&menus.Subject25, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject25.Text)
	})

	bot.Handle(&menus.Subject26, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject26.Text)
	})

	bot.Handle(&menus.Subject27, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject27.Text)
	})

	bot.Handle(&menus.Subject28, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject28.Text)
	})

	bot.Handle(&menus.Subject29, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject29.Text)
	})

	bot.Handle(&menus.Subject30, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject30.Text)
	})

	bot.Handle(&menus.Subject31, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject31.Text)
	})

	bot.Handle(&menus.Subject32, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject32.Text)
	})

	bot.Handle(&menus.Subject33, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject33.Text)
	})

	bot.Handle(&menus.Subject34, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, rdb, ctx, menus.Subject34.Text)
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

func handlingTriggersOnMessages(bot *telegram.Bot, rdb *redis.Client) {
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
}
