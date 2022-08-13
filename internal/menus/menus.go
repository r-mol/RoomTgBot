package menus

import (
	"RoomTgBot/internal/commands"

	telegram "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	MainMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	RoomMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	AquaManMenu  = &telegram.ReplyMarkup{ResizeKeyboard: true}
	CleanManMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

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

	// Clean man menu buttons.
	BtnCleanRoom = CleanManMenu.Text(commands.CommandCleanRoom)

	BtnBack = MainMenu.Text(commands.CommandBack)
)

func InitializeMenus() {
	MainMenu.Reply(
		MainMenu.Row(BtnRoom, BtnNews),
		MainMenu.Row(BtnExam, BtnSettings),
	)

	RoomMenu.Reply(
		RoomMenu.Row(BtnShop, BtnAquaMan, BtnCleanMan),
		MainMenu.Row(BtnBack),
	)

	AquaManMenu.Reply(
		AquaManMenu.Row(BtnBringWater),
		MainMenu.Row(BtnBack),
	)

	CleanManMenu.Reply(
		CleanManMenu.Row(BtnCleanRoom),
		MainMenu.Row(BtnBack),
	)
}

func GetMenus() map[string]*telegram.ReplyMarkup {
	allMenus := map[string]*telegram.ReplyMarkup{}

	allMenus[commands.CommandStart] = MainMenu
	allMenus[commands.CommandRoom] = RoomMenu
	allMenus[commands.CommandAquaMan] = AquaManMenu

	return allMenus
}
