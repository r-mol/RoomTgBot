package exam

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"

	"context"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

func GetSetExam(contex context.Context, rdb *redis.Client, ctx telegram.Context, subjectName string) error {
	curState, err := state.GetCurStateFromRDB(contex, rdb, ctx)
	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	if curState.StateName == commands.CommandExamDone {
		files := curState.Files
		photos := curState.Photos
		err = setExam(subjectName, files, photos)

		if err != nil {
			return err
		}

		commandFrom := curState.StateName
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commands.CommandStart)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful set...", menus.MainMenu)
	} else if curState.StateName == commands.CommandGetExam {
		files, photos, err := getExam(subjectName)
		if err != nil {
			return err
		}

		curState.Files = files
		curState.Photos = photos
		err = curState.SendAllAvailableMessages(ctx, nil)
		if err != nil {
			return err
		}

		curState.RemoveAll()
		commandFrom := curState.StateName
		err = state.CheckOfUserState(contex, rdb, ctx, commandFrom, commands.CommandStart)
		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return ctx.Send("Exams successful get...", menus.MainMenu)
	}

	return ctx.Send("Please restart bot ✨")
}

func setExam(subjectName string, files []*telegram.Document, photos []*telegram.Photo) error {
	return nil
}

func getExam(subjectName string) ([]*telegram.Document, []*telegram.Photo, error) {
	return []*telegram.Document{}, []*telegram.Photo{}, nil
}
