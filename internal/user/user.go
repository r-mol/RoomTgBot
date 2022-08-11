package user

import (
	telegram "gopkg.in/telebot.v3"
)

type User struct {
	ID int64 `json:"id"`

	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
	CurState  string `json:"cur_state"`
}

func CreateUser(bot *telegram.Bot, ctx telegram.Context, newUser *User) error {
	*newUser = User{
		ID:        ctx.Sender().ID,
		FirstName: ctx.Sender().FirstName,
		Username:  ctx.Sender().Username,
		IsBot:     ctx.Sender().IsBot,
		CurState:  " ",
	}

	if newUser.IsBot {
		defer bot.Stop()
		return ctx.Send("You are fucking bot...")
	}

	return nil
}
