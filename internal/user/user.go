package user

import (
	"strconv"

	telegram "gopkg.in/telebot.v3"
)

type User struct {
	ID int64 `json:"id"`

	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

func CreateUser(bot *telegram.Bot, ctx telegram.Context, newUser *User) error {
	id := ctx.Sender().ID

	// TODO if user with id exist in database then:
	//   newUser = database(id)
	//   else create new user

	*newUser = User{
		ID: id,

		FirstName: ctx.Sender().FirstName,
		Username:  ctx.Sender().Username,
		IsBot:     ctx.Sender().IsBot,
	}

	if newUser.IsBot {
		defer bot.Stop()
		return ctx.Send("You are fucking bot...")
	}

	return nil
}

func (u *User) Recipient() string {
	return strconv.FormatInt(u.ID, 10)
}
