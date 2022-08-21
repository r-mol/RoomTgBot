package settings

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/state"

	"context"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

func ChangeWantedNotificationsOf(contex context.Context, rdb *redis.Client, ctx telegram.Context, notificationKey string) error {
	states := state.States{}

	err := state.GetStatesFromRDB(contex, rdb, ctx.Sender().ID, &states)
	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	notificationState := states[consts.Notification]
	waitedNotification := notificationState.WaitedNotification

	if _, ok := waitedNotification[notificationKey]; ok {
		delete(waitedNotification, notificationKey)

		err = ctx.Send("This notification has been turn off.", menus.SettingsBackMenu)
		if err != nil {
			return nil
		}
	} else {
		waitedNotification[notificationKey] = struct{}{}
		err = ctx.Send("This notification has been turn on.", menus.SettingsBackMenu)
		if err != nil {
			return nil
		}
	}

	return notificationState.ChangeDataInState(contex, rdb, ctx.Sender().ID)
}
