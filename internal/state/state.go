package state

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/menus"
	"RoomTgBot/internal/user"
	"reflect"

	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

type State struct {
	StateName string `json:"state_name"`
	PrevState string `json:"prev_state"`
	Message
	Notifications `json:"map_notifications"`
	IsNow         bool `json:"is_now"`
}

type Notifications struct {
	WaitedNotification map[string]struct{} `json:"waited_notification"`
	Nfs                map[string]Messages `json:"nfs"`
}

type Message struct {
	Text   string              `json:"text"`
	Files  []telegram.Document `json:"document"`
	Photos []telegram.Photo    `json:"photo"`
}

type States map[string]*State
type Messages []Message

func CheckOfUserState(contex context.Context, rdb *redis.Client, ctx telegram.Context, prevCommand, initCommand string) error {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, ctx.Sender().ID, &states)
	if err != nil {
		return err
	}

	prevState, ok := states[prevCommand]
	if !ok {
		err = ctx.Send("Something bad happened,\nwe return you to the beginning...", menus.MainMenu)
		if err != nil {
			return err
		}

		err = resetToZeroState(contex, rdb, ctx.Sender().ID, states)
		if err != nil {
			return err
		}

		return fmt.Errorf("bad request prevState is not exist")
	}

	if !prevState.IsNow {
		err = ctx.Send("Something bad happened,\nwe return you to the beginning...", menus.MainMenu)
		if err != nil {
			return err
		}

		err = resetToZeroState(contex, rdb, ctx.Sender().ID, states)
		if err != nil {
			return err
		}

		return fmt.Errorf("bad request: prevState is not now")
	}

	prevState.IsNow = false

	curState, ok := states[initCommand]
	if !ok {
		curState = &State{
			StateName: initCommand,
			PrevState: prevState.StateName,
			Message: Message{
				Files:  []telegram.Document{},
				Photos: []telegram.Photo{},
			},
			IsNow: true,
		}
	} else {
		curState.IsNow = true
	}

	if prevState.PrevState != initCommand {
		prevState.MoveMessagesTo(curState)
	}

	states[initCommand] = curState
	states[consts.InitState] = curState

	return SetStatesToRDB(contex, rdb, ctx.Sender().ID, &states)
}

func GetCurStateFromRDB(contex context.Context, rdb *redis.Client, id int64) (*State, error) {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, id, &states)
	if err != nil {
		return nil, err
	}

	curState := states[consts.InitState]

	return curState, nil
}

func GetStatesFromRDB(contex context.Context, rdb *redis.Client, id int64, sts *States) error {
	idString := strconv.FormatInt(id, consts.BaseForConvertToInt)

	stateBytes, err := rdb.Get(contex, idString).Result()

	switch {
	case err != nil:
		return err
	default:
		err = json.Unmarshal([]byte(stateBytes), &sts)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetStatesToRDB(contex context.Context, rdb *redis.Client, id int64, sts *States) error {
	idString := strconv.FormatInt(id, consts.BaseForConvertToInt)

	stateBytes, err := json.Marshal(sts)
	if err != nil {
		return err
	}

	err = rdb.Set(contex, idString, stateBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func resetToZeroState(contex context.Context, rdb *redis.Client, id int64, states States) error {
	curState, err := GetCurStateFromRDB(contex, rdb, id)
	if err != nil {
		return err
	}

	curState.IsNow = false
	states[curState.StateName] = curState

	newCurState := states[consts.CommandStart]

	newCurState.IsNow = true

	states[consts.InitState] = newCurState
	states[newCurState.StateName] = newCurState

	err = SetStatesToRDB(contex, rdb, id, &states)
	if err != nil {
		return err
	}

	return nil
}

func GetSetOfAvailableChattingStates() map[string]struct{} {
	setOfStates := map[string]struct{}{}

	setOfStates[consts.CommandUploadNews] = struct{}{}
	setOfStates[consts.CommandUploadPurchase] = struct{}{}
	setOfStates[consts.CommandUploadExam] = struct{}{}

	return setOfStates
}

func GetMapOfWaitedNotifications() map[string]struct{} {
	waitedNotification := map[string]struct{}{}

	waitedNotification[consts.NotificationShop] = struct{}{}
	waitedNotification[consts.NotificationNews] = struct{}{}
	waitedNotification[consts.NotificationExam] = struct{}{}
	waitedNotification[consts.NotificationMoney] = struct{}{}
	waitedNotification[consts.NotificationCleaning] = struct{}{}
	waitedNotification[consts.CommandAquaManIN] = struct{}{}
	waitedNotification[consts.CommandCleanManIN] = struct{}{}

	return waitedNotification
}

func (state *State) MoveMessagesTo(curState *State) {
	if state.Text != "" {
		curState.Text = state.Text
		state.Text = ""
	}

	if len(state.Files) != 0 {
		curState.Files = append(curState.Files, state.Files...)
		state.Files = []telegram.Document{}
	}

	if len(state.Photos) != 0 {
		curState.Photos = append(curState.Photos, state.Photos...)
		state.Photos = []telegram.Photo{}
	}
}

func (state *State) RemoveAll() {
	state.Text = ""
	state.Files = []telegram.Document{}
	state.Photos = []telegram.Photo{}
}

func (state *State) SendAllAvailableMessage(bot *telegram.Bot, u *telegram.User, message Message, menu *telegram.ReplyMarkup) error {
	var err error

	if reflect.DeepEqual(message, Message{}) {
		message = state.Message
	}

	if message.Text == "" {
		message.Text = "Empty Text"
	}

	_, err = bot.Send(u, message.Text, menu)
	if err != nil {
		return err
	}

	if len(message.Photos) != 0 {
		for index := range message.Photos {
			_, err = bot.Send(u, &message.Photos[index])
			if err != nil {
				return err
			}
		}
	}

	if len(message.Files) != 0 {
		for index := range message.Files {
			_, err = bot.Send(u, &message.Files[index])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (state *State) ChangeDataInState(contex context.Context, rdb *redis.Client, id int64) error {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, id, &states)
	if err != nil {
		return err
	}

	if states[consts.InitState].StateName == state.StateName {
		states[consts.InitState] = state
	}

	states[state.StateName] = state

	return SetStatesToRDB(contex, rdb, id, &states)
}

func ReturnToStartState(contex context.Context, rdb *redis.Client, ctx telegram.Context) error {
	state, err := GetCurStateFromRDB(contex, rdb, ctx.Sender().ID)
	if err != nil {
		return err
	}

	commandFrom := state.StateName

	err = CheckOfUserState(contex, rdb, ctx, commandFrom, consts.CommandStart)
	if err != nil {
		return err
	}

	return nil
}

func SetNotificationToAllUsers(contex context.Context, rdb *redis.Client, kindNotification string, message Message) error {
	allUsers := map[int64]telegram.User{}

	err := user.GetUsersFromDB(contex, rdb, allUsers)
	if err != nil {
		return err
	}

	for _, u := range allUsers {
		err := SetNotificationToUser(contex, rdb, u.ID, kindNotification, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetNotificationToUser(contex context.Context, rdb *redis.Client, id int64, keyOfNotification string, message Message) error {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, id, &states)

	if err != nil {
		return err
	}

	curState := states[consts.Notification]

	if curState.Nfs == nil {
		curState.Nfs = map[string]Messages{}
	}

	messages := curState.Nfs[keyOfNotification]

	if messages == nil {
		messages = []Message{}
	}

	messages = append(messages, message)

	curState.Nfs[keyOfNotification] = messages

	states[consts.Notification] = curState

	err = SetStatesToRDB(contex, rdb, id, &states)
	if err != nil {
		return err
	}

	return nil
}

func SendSpecialNotificationByKey(contex context.Context, bot *telegram.Bot, u *telegram.User, rdb *redis.Client, notificationKey string) error {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, u.ID, &states)
	if err != nil {
		return err
	}

	notificationState := states[consts.Notification]
	if len(notificationState.Nfs) == 0 {
		_, err = bot.Send(u, "Unfortunately, there are not smth new 🤷🏼", menus.MainMenu)
		return err
	}

	messages := notificationState.Nfs[notificationKey]

	if len(messages) == 0 {
		_, err = bot.Send(u, "Unfortunately, there are not smth new 🤷🏼", menus.MainMenu)
		return err
	}

	allMenus := menus.GetMenus()

	value, ok := allMenus[notificationKey]
	if !ok {
		value = allMenus[states[consts.InitState].StateName]
	}

	for _, message := range messages {
		err = notificationState.SendAllAvailableMessage(bot, u, message, value)
		if err != nil {
			return err
		}
	}

	notificationState.Nfs[notificationKey] = Messages{}

	return SetStatesToRDB(contex, rdb, u.ID, &states)
}

func CheckUserOnAvailableNotifications(contex context.Context, bot *telegram.Bot, u *telegram.User, rdb *redis.Client) error {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, u.ID, &states)
	if err != nil {
		return err
	}

	notificationState := states[consts.Notification]
	if notificationState == nil {
		states[consts.Notification] = &State{
			StateName: consts.Notification,
			PrevState: consts.CommandStart,
			Notifications: Notifications{
				Nfs:                map[string]Messages{},
				WaitedNotification: map[string]struct{}{},
			},
		}

		err = SetStatesToRDB(contex, rdb, u.ID, &states)

		return err
	}

	if len(notificationState.Nfs) == 0 {
		return nil
	}

	mapNotifications := notificationState.Nfs
	waitedNotification := notificationState.WaitedNotification

	for key := range waitedNotification {
		messages, ok := mapNotifications[key]
		if !ok {
			continue
		}

		allMenus := menus.GetMenus()

		for _, message := range messages {
			err = notificationState.SendAllAvailableMessage(bot, u, message, allMenus[key])
			if err != nil {
				return err
			}
		}

		notificationState.Nfs[key] = Messages{}

		err = SetStatesToRDB(contex, rdb, u.ID, &states)
		if err != nil {
			return err
		}
	}

	return nil
}
