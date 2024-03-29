package bot

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/exam"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/mongodb"
	"RoomTgBot/internal/settings"
	"RoomTgBot/internal/state"
	"RoomTgBot/internal/user"
	"strconv"

	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	telegram "gopkg.in/telebot.v3"
)

var contex = context.Background()

func handling(bot *telegram.Bot, rdb *redis.Client, mdb *mongo.Client) {
	menus.InitializeMenus()

	allMenus := menus.GetMenus()

	handlingStart(bot, rdb, mdb)
	functionalHandling(bot, rdb, allMenus)
	handlingRoomMenu(bot, rdb, mdb)
	handlingNewsMenu(bot, rdb, allMenus)
	handlingExamMenu(bot, rdb)
	handlingSettingsMenu(bot, rdb)
	handlingTriggersOnMessages(bot, rdb)
}

func handlingStart(bot *telegram.Bot, rdb *redis.Client, mdb *mongo.Client) {
	bot.Handle(consts.CommandStart, func(ctx telegram.Context) error {
		var notificationState *state.State

		err := user.CreateUser(contex, rdb, mdb, bot, ctx)
		if err != nil {
			return err
		}

		log.Println("User is authorized")

		curState := &state.State{
			StateName: consts.CommandStart,
			IsNow:     true,
		}

		states := state.States{}

		err = state.GetStatesFromRDB(contex, rdb, ctx.Sender().ID, &states)
		switch err {
		case redis.Nil:
			waitedNotification := state.GetMapOfWaitedNotifications()

			notificationState = &state.State{StateName: consts.Notification,
				PrevState: consts.CommandStart,
				Notifications: state.Notifications{
					Nfs:                map[string]state.Messages{},
					WaitedNotification: waitedNotification,
				},
			}
		case nil:
			notificationState = states[consts.Notification]
		default:
			return err
		}

		states = state.States{}
		states[consts.Notification] = notificationState
		states[consts.CommandStart] = curState
		states[consts.InitState] = curState

		err = state.SetStatesToRDB(contex, rdb, ctx.Sender().ID, &states)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Nice to meet you "+ctx.Sender().FirstName+" !!!", menus.MainMenu)
	})
}

func functionalHandling(bot *telegram.Bot, rdb *redis.Client, allMenus map[string]*telegram.ReplyMarkup) {
	bot.Handle(&menus.BtnBack, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		if curState.StateName == consts.CommandStart {
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

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		return ctx.Send("Welcome back 🛑", allMenus[commandTo])
	})
}

func handlingRoomMenu(bot *telegram.Bot, rdb *redis.Client, mdb *mongo.Client) {
	bot.Handle(&menus.BtnRoom, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandStart, consts.CommandRoom)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the room menu...", menus.RoomMenu)
	})

	handlingDebter(bot, rdb, mdb)
	handlingShopMenu(bot, rdb)
	handlingAquaMan(bot, rdb)
	handlingCleanMan(bot, rdb)
}

func handlingDebter(bot *telegram.Bot, rdb *redis.Client, mdb *mongo.Client) {
	bot.Handle(&menus.BtnNotInInnoAQ, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}

		u := usersMap[ctx.Sender().ID]
		u.IsAbsent = true

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, u)
		if err != nil {
			return err
		}

		err = FindInitAquaMan()
		if err != nil {
			return err
		}

		return ctx.Send("Thanks for the answer, have a good time 😊", menus.MainMenu)
	})

	bot.Handle(&menus.BtnNotInInnoCR, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}

		u := usersMap[ctx.Sender().ID]
		u.IsAbsent = true

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, u)
		if err != nil {
			return err
		}

		err = FindInitCleanMan()
		if err != nil {
			return err
		}

		return ctx.Send("Thanks for the answer, have a good time 😊", menus.MainMenu)
	})

	bot.Handle(&menus.BtnCantAQ, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = FindInitAquaMan()
		if err != nil {
			return err
		}

		return ctx.Send("We heard you, please don't let this happen again 🥺", menus.MainMenu)
	})

	bot.Handle(&menus.BtnCantCR, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = FindInitCleanMan()
		if err != nil {
			return err
		}

		return ctx.Send("We heard you, please don't let this happen again 🥺", menus.MainMenu)
	})

	bot.Handle(&menus.BtnAquaManIN, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}

		usersMap, err = user.IncreaseScore(ctx.Sender().ID, usersMap, consts.CommandAquaManIN, consts.InitialActivityList)
		if err != nil {
			return err
		}

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, usersMap[ctx.Sender().ID])
		if err != nil {
			return err
		}

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(&menus.BtnCleanManIN, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}

		usersMap, err = user.IncreaseScore(ctx.Sender().ID, usersMap, consts.CommandCleanManIN, consts.InitialActivityList)
		if err != nil {
			return err
		}

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, usersMap[ctx.Sender().ID])
		if err != nil {
			return err
		}

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})
}

func handlingShopMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnShop, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandRoom, consts.CommandShop)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the shop menu...", menus.ShopMenu)
	})

	bot.Handle(&menus.BtnUploadPurchase, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandShop, consts.CommandUploadPurchase)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send:\n1. Name of purchase\n2. Image of purchase\n3. Cost of purchase", menus.ShopUploadMenu)
	})

	bot.Handle(&menus.BtnPurchaseDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandUploadPurchase, consts.CommandPurchaseDone)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Please, check your draft of purchase:")
		if err != nil {
			return err
		}

		curState.Message.Text = "🛍 Purchase report 🛍\n\n" + curState.Message.Text

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return curState.SendAllAvailableMessage(bot, ctx.Sender(), state.Message{}, menus.PostPurchaseMenu)
	})

	bot.Handle(&menus.BtnPostPurchase, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = state.SetNotificationToAllUsers(contex, rdb, consts.NotificationShop, curState.Message)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		err = state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Your purchase has been send to database 📨", menus.MainMenu)
	})

	bot.Handle(&menus.BtnCheckShopping, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Here you can check new purchases:")
		if err != nil {
			return err
		}

		return state.SendSpecialNotificationByKey(contex, bot, ctx.Sender(), rdb, consts.NotificationShop)
	})
}

func handlingAquaMan(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnAquaMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandRoom, consts.CommandAquaMan)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are aqua-man...", menus.AquaManMenu)
	})

	bot.Handle(consts.CommandBringWater, func(ctx telegram.Context) error {
		fmt.Println("flag0")
		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}
		fmt.Println("flag1")

		usersMap, err = user.IncreaseScore(ctx.Sender().ID, usersMap, consts.CommandAquaManIN, consts.InitialActivityList)
		if err != nil {
			return err
		}
		fmt.Println("flag2")

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, usersMap[ctx.Sender().ID])
		if err != nil {
			return err
		}
		fmt.Println("flag3")

		err = state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}
		fmt.Println("flag4")

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})

	bot.Handle(consts.CommandWaterIsOver, func(ctx telegram.Context) error {
		err := FindInitAquaMan()
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})
	bot.Handle(consts.CommandGetScore, func(ctx telegram.Context) error {
		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}
		answ := "Total score is:"
		for _, user := range usersMap {
			answ += "\n" + user.TelegramUsername + ":"
			for _, val := range user.ScoreList {
				answ += " " + strconv.Itoa(val)
			}
		}
		return ctx.Send(answ)
	})

}

var prevIDAQ = int64(0)

func FindInitAquaMan() error {
	usersMap, err := user.MongoGetMap(contex, mdb)
	if err != nil {
		return err
	}

	prevIDAQ, err = user.NextInOrder(prevIDAQ, usersMap, consts.InitialActivityList[consts.CommandAquaManIN].MongoID)
	if err != nil {
		return err
	}

	message := state.Message{Text: "Please, bring the water to room."}

	err = state.SetNotificationToUser(contex, rdb, prevIDAQ, consts.CommandAquaManIN, message)
	if err != nil {
		return err
	}

	return nil
}

func handlingCleanMan(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnCleanMan, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandRoom, consts.CommandCleanMan)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are clean-man...", menus.CleanManMenu)
	})

	bot.Handle(consts.CommandCleanRoom, func(ctx telegram.Context) error {
		usersMap, err := user.MongoGetMap(contex, mdb)
		if err != nil {
			return err
		}

		usersMap, err = user.IncreaseScore(ctx.Sender().ID, usersMap, consts.CommandCleanManIN, consts.InitialActivityList)
		if err != nil {
			return err
		}

		err = mongodb.UpdateOne(contex, mdb, consts.MongoUsersCollection, usersMap[ctx.Sender().ID])
		if err != nil {
			return err
		}

		err = state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("We really appreciate your contribution in maintaining the room 💪🏽", menus.MainMenu)
	})
}

var prevIDC = int64(0)

func FindInitCleanMan() error {
	usersMap, err := user.MongoGetMap(contex, mdb)
	if err != nil {
		return err
	}

	prevIDC, err = user.NextInOrder(prevIDC, usersMap, consts.InitialActivityList[consts.CommandCleanManIN].MongoID)
	if err != nil {
		return err
	}

	message := state.Message{Text: "Please, clean room."}

	err = state.SetNotificationToUser(contex, rdb, prevIDC, consts.CommandCleanManIN, message)
	if err != nil {
		return err
	}

	return nil
}

func NotifyAboutCleaning() error {
	message := state.Message{Text: "Please, don't forget that tomorrow will be cleaning."}

	err := state.SetNotificationToAllUsers(contex, rdb, consts.NotificationCleaning, message)
	if err != nil {
		return err
	}

	return nil
}

func NotifyAboutMoney() error {
	message := state.Message{Text: "Please, don't forget to pay 100rub to room account ."}

	err := state.SetNotificationToAllUsers(contex, rdb, consts.NotificationMoney, message)
	if err != nil {
		return err
	}

	return nil
}

func PutNotAbsentToAllUsers() error {
	return user.NotAbsentAllUsers(contex, mdb)
}

func handlingNewsMenu(bot *telegram.Bot, rdb *redis.Client, allMenus map[string]*telegram.ReplyMarkup) {
	bot.Handle(&menus.BtnNews, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandStart, consts.CommandNews)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the news menu...", menus.NewsMenu)
	})

	bot.Handle(&menus.BtnUploadNews, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandNews, consts.CommandUploadNews)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send text/files/photos messages to create news:", menus.NewsUploadMenu)
	})

	bot.Handle(&menus.BtnNewsDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandUploadNews, consts.CommandNewsDone)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Please, check your draft of news:")
		if err != nil {
			return err
		}

		curState.Message.Text = "‼️‼️ News ‼️‼️\n\n" + curState.Message.Text

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return curState.SendAllAvailableMessage(bot, ctx.Sender(), state.Message{}, menus.PostNewsMenu)
	})

	bot.Handle(&menus.BtnPostNews, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = state.SetNotificationToAllUsers(contex, rdb, consts.NotificationNews, curState.Message)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		err = state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Your news has been posted 📨", menus.MainMenu)
	})

	bot.Handle(&menus.BtnCheckNews, func(ctx telegram.Context) error {
		err := state.ReturnToStartState(contex, rdb, ctx)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Here you can check News:")
		if err != nil {
			return err
		}

		return state.SendSpecialNotificationByKey(contex, bot, ctx.Sender(), rdb, consts.NotificationNews)
	})

	bot.Handle(&menus.BtnDeleteDraft, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		if curState.StateName == consts.CommandStart {
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

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		return ctx.Send("All you messages have been removed", allMenus[commandTo])
	})
}

func handlingExamMenu(bot *telegram.Bot, rdb *redis.Client) {
	handlingSubjects(bot, rdb)

	bot.Handle(&menus.BtnExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandStart, consts.CommandExam)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the exam menu...", menus.ExamMenu)
	})

	bot.Handle(&menus.BtnUploadExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandExam, consts.CommandUploadExam)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, send files:", menus.ExamUploadMenu)
	})

	bot.Handle(&menus.BtnGetExam, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandExam, consts.CommandGetExam)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Please, choose subject from list:", menus.SubjectMenu)
	})

	bot.Handle(&menus.BtnExamDone, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandUploadExam, consts.CommandExamDone)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessage(bot, ctx.Sender(), state.Message{}, menus.MainMenu)
		if err != nil {
			return err
		}

		return ctx.Send("Please, check your files of exam:", menus.SubjectMenu)
	})
}

func handlingSubjects(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.Subject1, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject1.Text)
	})

	bot.Handle(&menus.Subject2, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject2.Text)
	})

	bot.Handle(&menus.Subject3, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject3.Text)
	})

	bot.Handle(&menus.Subject4, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject4.Text)
	})

	bot.Handle(&menus.Subject5, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject5.Text)
	})

	bot.Handle(&menus.Subject6, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject6.Text)
	})

	bot.Handle(&menus.Subject7, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject7.Text)
	})

	bot.Handle(&menus.Subject8, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject8.Text)
	})

	bot.Handle(&menus.Subject9, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject9.Text)
	})

	bot.Handle(&menus.Subject10, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject10.Text)
	})

	bot.Handle(&menus.Subject11, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject11.Text)
	})

	bot.Handle(&menus.Subject12, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject12.Text)
	})

	bot.Handle(&menus.Subject13, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject13.Text)
	})

	bot.Handle(&menus.Subject14, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject14.Text)
	})

	bot.Handle(&menus.Subject15, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject15.Text)
	})

	bot.Handle(&menus.Subject16, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject16.Text)
	})

	bot.Handle(&menus.Subject17, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject17.Text)
	})

	bot.Handle(&menus.Subject18, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject18.Text)
	})

	bot.Handle(&menus.Subject19, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject19.Text)
	})

	bot.Handle(&menus.Subject20, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject20.Text)
	})

	bot.Handle(&menus.Subject21, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject21.Text)
	})

	bot.Handle(&menus.Subject22, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject22.Text)
	})

	bot.Handle(&menus.Subject23, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject23.Text)
	})

	bot.Handle(&menus.Subject24, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject24.Text)
	})

	bot.Handle(&menus.Subject25, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject25.Text)
	})

	bot.Handle(&menus.Subject26, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject26.Text)
	})

	bot.Handle(&menus.Subject27, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject27.Text)
	})

	bot.Handle(&menus.Subject28, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject28.Text)
	})

	bot.Handle(&menus.Subject29, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject29.Text)
	})

	bot.Handle(&menus.Subject30, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject30.Text)
	})

	bot.Handle(&menus.Subject31, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject31.Text)
	})

	bot.Handle(&menus.Subject32, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject32.Text)
	})

	bot.Handle(&menus.Subject33, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject33.Text)
	})

	bot.Handle(&menus.Subject34, func(ctx telegram.Context) error {
		return exam.GetSetExam(contex, bot, mdb, rdb, ctx, menus.Subject34.Text)
	})
}

func handlingSettingsMenu(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(&menus.BtnSettings, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandStart, consts.CommandSettings)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Now you are in the settings menu...", menus.SettingsMenu)
	})

	bot.Handle(&menus.BtnNotificationSettings, func(ctx telegram.Context) error {
		err := state.CheckOfUserState(contex, rdb, ctx, consts.CommandSettings, consts.CommandNotificationSettings)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		err = ctx.Send("Press on button to choose ", menus.SettingsBackMenu)
		if err != nil {
			return err
		}

		return ctx.Send("Which notification turn on/off:", menus.SettingsOfNotifications)
	})

	bot.Handle(&menus.BtnShopNotification, func(ctx telegram.Context) error {
		return settings.ChangeWantedNotificationsOf(contex, rdb, ctx, menus.BtnShopNotification.Text)
	})

	bot.Handle(&menus.BtnNewsNotification, func(ctx telegram.Context) error {
		return settings.ChangeWantedNotificationsOf(contex, rdb, ctx, menus.BtnNewsNotification.Text)
	})

	bot.Handle(&menus.BtnExamNotification, func(ctx telegram.Context) error {
		return settings.ChangeWantedNotificationsOf(contex, rdb, ctx, menus.BtnExamNotification.Text)
	})

	bot.Handle(&menus.BtnMoneyNotification, func(ctx telegram.Context) error {
		return settings.ChangeWantedNotificationsOf(contex, rdb, ctx, menus.BtnMoneyNotification.Text)
	})

	bot.Handle(&menus.BtnCleaningNotification, func(ctx telegram.Context) error {
		return settings.ChangeWantedNotificationsOf(contex, rdb, ctx, menus.BtnCleaningNotification.Text)
	})
}

func handlingTriggersOnMessages(bot *telegram.Bot, rdb *redis.Client) {
	bot.Handle(telegram.OnText, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
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

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(telegram.OnDocument, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableChattingStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("You can not write here or you send unavailable command...")
		}

		curState.Files = append(curState.Files, *ctx.Message().Document)

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Handle(telegram.OnPhoto, func(ctx telegram.Context) error {
		curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		setOfStates := state.GetSetOfAvailableChattingStates()

		if _, ok := setOfStates[curState.StateName]; !ok {
			return ctx.Send("You can not write here or you send unavailable command...")
		}

		curState.Photos = append(curState.Photos, *ctx.Message().Photo)

		err = curState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
		if err != nil {
			return err
		}

		return nil
	})
}
