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

type Message struct {
	Text   string              `json:"text"`
	Files  []telegram.Document `json:"document"`
	Photos []telegram.Photo    `json:"photo"`
}

type Messages struct {
	Index        int       `json:"index"`
	ListMessages []Message `json:"list_messages"`
}

type Notifications struct {
	MapNotifications map[string]Messages `json:"map_notifications"`
}

type State struct {
	StateName string `json:"init_state"`
	PrevState string `json:"prev_state"`
	Message
	Messages
	Notifications
	IsNow bool `json:"is_now"`
}

type States map[string]*State

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

	if err != nil {
		return err
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

	if err != nil {
		return err
	}

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

	err = SetStatesToRDB(contex, rdb, ctx.Sender().ID, &states)

	if err != nil {
		return err
	}

	return err
}

func GetCurStateFromRDB(contex context.Context, rdb *redis.Client, ID int64) (*State, error) {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, ID, &states)

	if err != nil {
		return nil, err
	}

	curState := states[consts.InitState]

	return curState, nil
}

func GetStatesFromRDB(contex context.Context, rdb *redis.Client, ID int64, sts *States) error {
	id := strconv.FormatInt(ID, consts.BaseForConvertToInt)

	stateBytes, err := rdb.Get(contex, id).Result()

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

func SetStatesToRDB(contex context.Context, rdb *redis.Client, ID int64, sts *States) error {
	id := strconv.FormatInt(ID, consts.BaseForConvertToInt)

	stateBytes, err := json.Marshal(sts)

	if err != nil {
		return err
	}

	err = rdb.Set(contex, id, stateBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func resetToZeroState(contex context.Context, rdb *redis.Client, ID int64, states States) error {
	curState, err := GetCurStateFromRDB(contex, rdb, ID)

	if err != nil {
		return err
	}

	curState.IsNow = false
	states[curState.StateName] = curState

	newCurState := states[consts.CommandStart]

	newCurState.IsNow = true

	states[consts.InitState] = newCurState
	states[newCurState.StateName] = newCurState

	err = SetStatesToRDB(contex, rdb, ID, &states)

	if err != nil {
		return err
	}

	return nil
}

func GetSetOfAvailableChattingStates() map[string]struct{} {
	setOfStates := map[string]struct{}{}

	setOfStates[consts.CommandNews] = struct{}{}
	setOfStates[consts.CommandUploadPurchase] = struct{}{}
	setOfStates[consts.CommandUploadExam] = struct{}{}

	return setOfStates
}

func GetSetOfAvailableListStates() map[string]struct{} {
	setOfStates := map[string]struct{}{}

	setOfStates[consts.CommandCheck] = struct{}{}
	setOfStates[consts.Notification] = struct{}{}

	return setOfStates
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
	state.ListMessages = []Message{}
}

func (state *State) SendAllAvailableMessage(bot *telegram.Bot, user *telegram.User, message Message, menu *telegram.ReplyMarkup) error {
	var err error

	if reflect.DeepEqual(message, Message{}) {
		message = state.Message
	}

	if len(message.Files) != 0 {
		for _, file := range message.Files {
			_, err = bot.Send(user, file)
			if err != nil {
				return err
			}
		}
	}

	if len(message.Photos) != 0 {
		for _, photo := range message.Photos {
			_, err = bot.Send(user, photo)
			if err != nil {
				return err
			}
		}
	}

	if message.Text == "" {
		message.Text = "Empty Text"
	}

	_, err = bot.Send(user, message.Text, menu)

	return err
}

func (state *State) SendAllAvailableMessages(bot *telegram.Bot, user *telegram.User, message Message, menu *telegram.ReplyMarkup) error {
	var err error

	if reflect.DeepEqual(message, Message{}) {
		message = state.Message
	}

	if len(message.Files) != 0 {
		for _, file := range message.Files {
			_, err = bot.Send(user, file)
			if err != nil {
				return err
			}
		}
	}

	if len(message.Photos) != 0 {
		for _, photo := range message.Photos {
			_, err = bot.Send(user, photo)
			if err != nil {
				return err
			}
		}
	}

	if message.Text != "" {
		_, err = bot.Send(user, message.Text)
		if err != nil {
			return err
		}
	}

	_, err = bot.Send(user, "You can list in the items:", menu)

	return err
}

func (state *State) ChangeDataInState(contex context.Context, rdb *redis.Client, ID int64) error {
	states := States{}

	err := GetStatesFromRDB(contex, rdb, ID, &states)

	if err != nil {
		return err
	}

	if states[consts.InitState].StateName == state.StateName {
		states[consts.InitState] = state
	}

	states[state.StateName] = state

	err = SetStatesToRDB(contex, rdb, ID, &states)

	if err != nil {
		return err
	}

	return nil
}

func (state *State) GetNextMessageOfList() Message {
	if len(state.ListMessages)-1 == state.Index {
		state.Index = 0
	} else {
		state.Index++
	}

	return state.ListMessages[state.Index]
}

func (state *State) GetPrevMessageOfList() Message {
	if state.Index == 0 {
		state.Index = len(state.ListMessages) - 1
	} else {
		state.Index--
	}

	return state.ListMessages[state.Index]
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
	// TODO Get all users from database
	allUsers := []user.User{}

	for _, u := range allUsers {
		err := SetNotificationToUser(contex, rdb, u.ID, kindNotification, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetNotificationToUser(contex context.Context, rdb *redis.Client, ID int64, keyOfNotification string, message Message) error {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, ID, &states)

	if err != nil {
		return err
	}

	curState := states[consts.Notification]

	if curState.MapNotifications == nil {
		curState.MapNotifications = map[string]Messages{}
	}

	listMessages := curState.MapNotifications[keyOfNotification].ListMessages
	if listMessages == nil {
		listMessages = []Message{}
	}

	listMessages = append(listMessages, message)

	curState.MapNotifications[keyOfNotification] = Messages{ListMessages: listMessages}

	states[consts.Notification] = curState

	err = SetStatesToRDB(contex, rdb, ID, &states)

	if err != nil {
		return err
	}

	return nil
}

func SendSpetialNotificationByKey(bot *telegram.Bot, user *telegram.User, contex context.Context, rdb *redis.Client, key string) error {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, user.ID, &states)

	if err != nil {
		return err
	}

	notificationState := states[consts.Notification]
	if len(notificationState.MapNotifications) == 0 {
		return nil
	}

	messages := notificationState.MapNotifications[key]

	allMenus := menus.GetMenus()

	for _, message := range messages.ListMessages {
		err = notificationState.SendAllAvailableMessages(bot, user, message, allMenus[key])
		if err != nil {
			return err
		}
	}

	notificationState.MapNotifications[key] = Messages{}

	return nil
}

func CheckUserOnAvaliableNotifications(bot *telegram.Bot, user *telegram.User, contex context.Context, rdb *redis.Client) error {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, user.ID, &states)

	if err != nil {
		return err
	}

	notificationState := states[consts.Notification]
	if notificationState == nil {
		states[consts.Notification] = &State{StateName: consts.Notification, PrevState: consts.CommandStart}
		err = SetStatesToRDB(contex, rdb, user.ID, &states)
		return err
	} else if len(notificationState.MapNotifications) == 0 {
		return nil
	}

	mapNotifications := notificationState.MapNotifications

	// TODO get user preferences of notification from data base
	waitedNotification := map[string]struct{}{}
	waitedNotification[consts.CommandAquaManIN] = struct{}{}

	for key := range waitedNotification {
		messages, ok := mapNotifications[key]
		if !ok {
			continue
		}

		allMenus := menus.GetMenus()

		for _, message := range messages.ListMessages {
			err = notificationState.SendAllAvailableMessage(bot, user, message, allMenus[key])
			if err != nil {
				return err
			}
		}

		notificationState.MapNotifications[key] = Messages{}
		err = SetStatesToRDB(contex, rdb, user.ID, &states)
		if err != nil {
			return err
		}
	}

	return nil
}
