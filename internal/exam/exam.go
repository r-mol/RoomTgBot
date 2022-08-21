package exam

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"

	"context"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

func GetSetExam(contex context.Context, bot *telegram.Bot, rdb *redis.Client, ctx telegram.Context, subjectName string) error {
	curState, err := state.GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	if curState.StateName == consts.CommandExamDone {
		files := curState.Files
		photos := curState.Photos

		err = setExam(subjectName, files, photos)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		commandFrom := curState.StateName

		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, consts.CommandStart)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful set...", menus.MainMenu)
	} else if curState.StateName == consts.CommandGetExam {
		curState.Files, curState.Photos, err = getExam(subjectName)
		if err != nil {
			return err
		}

		err = curState.SendAllAvailableMessage(bot, ctx.Sender(), state.Message{}, menus.MainMenu)
		if err != nil {
			return err
		}

		curState.RemoveAll()

		commandFrom := curState.StateName

		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, consts.CommandStart)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful get...", menus.MainMenu)
	}

	return ctx.Send("Please restart bot ✨")
}

func setExam(subjectName string, files []telegram.Document, photos []telegram.Photo) error {
	return nil
}

func getExam(subjectName string) ([]telegram.Document, []telegram.Photo, error) {
	return []telegram.Document{}, []telegram.Photo{}, nil
}
