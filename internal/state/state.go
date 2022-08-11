package state

import (
	"RoomTgBot/internal/menus"

	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
)

type State struct {
	InitState string    `json:"init_state"`
	PrevState string    `json:"prev_state"`
	Message   string    `json:"message"`
	File      []os.File `json:"file"`
	IsNow     bool      `json:"is_now"`
}

func CheckOfUserState(contex context.Context, rdb *redis.Client, ctx telegram.Context, prevCommand, initCommand string) (*State, error) {
	prevState := &State{}
	curState := &State{}

	err := GetStateFromRDB(contex, rdb, prevState, prevCommand)

	log.Println(prevState)

	if err != nil {
		return nil, err
	}

	if !prevState.IsNow {
		return nil, ctx.Send("Something bad happened, we return you to the beginning", menus.MainMenu)
	}

	prevState.IsNow = false
	err = SetStateToRDB(contex, rdb, prevState)

	if err != nil {
		return nil, err
	}

	err = GetStateFromRDB(contex, rdb, curState, initCommand)

	switch err {
	case redis.Nil:
		curState = &State{
			InitState: initCommand,
			PrevState: prevState.InitState,
			IsNow:     true,
		}

		err = SetStateToRDB(contex, rdb, curState)

		if err != nil {
			return nil, err
		}
	default:
		curState.IsNow = true
		err = SetStateToRDB(contex, rdb, curState)

		if err != nil {
			return nil, err
		}

		log.Println(prevState)
		log.Println(curState)

		return nil, err
	}

	log.Println(prevState)
	log.Println(curState)

	return curState, nil
}

func GetStateFromRDB(contex context.Context, rdb *redis.Client, st *State, command string) error {
	stBytes, err := rdb.Get(contex, command).Result()

	switch {
	case err == redis.Nil:
		return err
	case err != nil:
		return err
	default:
		err = json.Unmarshal([]byte(stBytes), &st)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetStateToRDB(contex context.Context, rdb *redis.Client, st *State) error {
	stateBytes, err := json.Marshal(st)

	if err != nil {
		return err
	}

	err = rdb.Set(contex, st.InitState, stateBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
