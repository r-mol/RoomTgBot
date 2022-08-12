package state

import (
	"RoomTgBot/internal/commands"
	"RoomTgBot/internal/menus"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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

type States map[string]*State

const InitState = "init_state"

func CheckOfUserState(contex context.Context, rdb *redis.Client, ctx telegram.Context, prevCommand, initCommand string) error {
	states := States{}
	curState := &State{}

	err := GetStatesFromRDB(contex, rdb, ctx, &states)

	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	prevState, ok := states[prevCommand]
	if !ok {
		err = ctx.Send("Something bad happened, we return you to the beginning", menus.MainMenu)

		if err != nil {
			return err
		}

		err = resetToZeroState(contex, rdb, ctx)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return fmt.Errorf("bad request prevState is not exist")
	}

	log.Println(prevState)

	if err != nil {
		return err
	}

	if !prevState.IsNow {
		err = ctx.Send("Something bad happened, we return you to the beginning", menus.MainMenu)

		if err != nil {
			return err
		}

		err = resetToZeroState(contex, rdb, ctx)

		if err == redis.Nil {
			return ctx.Send("Please restart bot ✨")
		} else if err != nil {
			return err
		}

		return fmt.Errorf("bad request: prevState is not now")
	}

	prevState.IsNow = false
	states[prevCommand] = prevState

	if err != nil {
		return err
	}

	curState, ok = states[initCommand]

	if !ok {
		curState = &State{
			InitState: initCommand,
			PrevState: prevState.InitState,
			IsNow:     true,
		}
	} else {
		curState.IsNow = true
	}

	states[initCommand] = curState
	states[InitState] = curState

	err = SetStatesToRDB(contex, rdb, ctx, &states)

	if err != nil {
		return err
	}

	log.Println(prevState)
	log.Println(curState)

	return err
}

func GetCurStateFromRDB(contex context.Context, rdb *redis.Client, ctx telegram.Context) (*State, error) {
	states := States{}
	err := GetStatesFromRDB(contex, rdb, ctx, &states)

	if err == redis.Nil {
		return nil, ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return nil, err
	}

	curState := states[InitState]

	return curState, nil
}

func GetStatesFromRDB(contex context.Context, rdb *redis.Client, ctx telegram.Context, sts *States) error {
	id := strconv.FormatInt(ctx.Sender().ID, 10)

	stateBytes, err := rdb.Get(contex, id).Result()

	switch {
	case err == redis.Nil:
		sts = &States{}
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

func SetStatesToRDB(contex context.Context, rdb *redis.Client, ctx telegram.Context, sts *States) error {
	id := strconv.FormatInt(ctx.Sender().ID, 10)

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

func resetToZeroState(contex context.Context, rdb *redis.Client, ctx telegram.Context) error {
	states := States{}
	curState, err := GetCurStateFromRDB(contex, rdb, ctx)

	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	curState.IsNow = false

	states[curState.InitState] = curState

	newCurState := states[commands.CommandStart]

	newCurState.IsNow = true

	states[commands.CommandStart] = newCurState
	states[InitState] = newCurState

	err = SetStatesToRDB(contex, rdb, ctx, &states)

	if err == redis.Nil {
		return ctx.Send("Please restart bot ✨")
	} else if err != nil {
		return err
	}

	return nil
}
