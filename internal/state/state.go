package state

import (
	"RoomTgBot/internal/menus"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	telegram "gopkg.in/telebot.v3"
	"log"
	"os"
)

type State struct {
	InitState string    `json:"init_state"`
	PrevState string    `json:"prev_state"`
	Message   string    `json:"message"`
	File      []os.File `json:"file"`
	IsNow     bool      `json:"is_now"`
}

func CheckOfUserState(rdb *redis.Client, ctx telegram.Context, contex context.Context, prevCommand string, initCommand string) error {
	prevState := &State{}
	st := &State{}

	err := GetStateFromRDB(rdb, contex, prevState, prevCommand)

	log.Println(prevState)

	if err != nil {
		return err
	}

	if !prevState.IsNow {
		return ctx.Send("Something bad happened, we return you to the beginning", menus.MainMenu)
	}

	prevState.IsNow = false
	err = SetStateToRDB(rdb, contex, prevState)

	if err != nil {
		return err
	}

	err = GetStateFromRDB(rdb, contex, st, initCommand)

	switch err {
	case redis.Nil:
		st = &State{
			InitState: initCommand,
			PrevState: prevState.InitState,
			IsNow:     true,
		}

		err = SetStateToRDB(rdb, contex, st)

		if err != nil {
			return err
		}
	default:
		st.IsNow = true
		err = SetStateToRDB(rdb, contex, st)

		if err != nil {
			return err
		}
		log.Println(prevState)
		log.Println(st)
		return err
	}
	log.Println(prevState)
	log.Println(st)

	return nil
}

func GetStateFromRDB(rdb *redis.Client, contex context.Context, st *State, command string) error {
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

func SetStateToRDB(rdb *redis.Client, contex context.Context, st *State) error {
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
