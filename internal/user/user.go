package user

import (
	"RoomTgBot/internal/state"
	telegram "gopkg.in/telebot.v3"
)

type User struct {
	ID int64 `json:"id"`

	FirstName string       `json:"first_name"`
	Username  string       `json:"username"`
	IsBot     bool         `json:"is_bot"`
	CurState  *state.State `json:"cur_state"`
}

func CreateUser(bot *telegram.Bot, ctx telegram.Context, newUser *User) error {
	id := ctx.Sender().ID

	//TODO if user with id exist in database then:
	// newUser = database(id)
	// else create new user

	*newUser = User{
		ID: id,

		FirstName: ctx.Sender().FirstName,
		Username:  ctx.Sender().Username,
		IsBot:     ctx.Sender().IsBot,
		CurState:  &state.State{},
	}

	if newUser.IsBot {
		defer bot.Stop()
		return ctx.Send("You are fucking bot...")
	}

	return nil
}
