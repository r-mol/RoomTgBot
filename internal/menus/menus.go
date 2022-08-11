package menus

import (
	"RoomTgBot/internal/commands"
	telegram "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	MainMenu    = &telegram.ReplyMarkup{ResizeKeyboard: true}
	RoomMenu    = &telegram.ReplyMarkup{ResizeKeyboard: true}
	AquaManMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	// Main menu buttons.
	BtnRoom     = MainMenu.Text(commands.CommandRoom)
	BtnNews     = MainMenu.Text(commands.CommandNews)
	BtnExam     = MainMenu.Text(commands.CommandExam)
	BtnSettings = MainMenu.Text(commands.CommandSettings)

	// Room menu buttons.
	BtnShop     = RoomMenu.Text(commands.CommandShop)
	BtnAquaMan  = RoomMenu.Text(commands.CommandAquaMan)
	BtnCleanMan = RoomMenu.Text(commands.CommandCleanMan)

	// Aqua man menu buttons.
	BtnBringWater = AquaManMenu.Text(commands.CommandBringWater)

	BtnBack = MainMenu.Text(commands.CommandBack)
)
