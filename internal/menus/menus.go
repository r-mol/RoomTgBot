package menus

import (
	"RoomTgBot/internal/commands"

	telegram "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	MainMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}
	RoomMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ExamMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}
	NewsMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}
	SettingsMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	AquaManMenu      = &telegram.ReplyMarkup{ResizeKeyboard: true}
	CleanManMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ShopMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}
	PostNewsMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ShopUploadMenu   = &telegram.ReplyMarkup{ResizeKeyboard: true}
	PostPurchaseMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ExamUploadMenu   = &telegram.ReplyMarkup{ResizeKeyboard: true}

	// Main menu buttons.
	BtnRoom     = MainMenu.Text(commands.CommandRoom)
	BtnNews     = MainMenu.Text(commands.CommandNews)
	BtnExam     = MainMenu.Text(commands.CommandExam)
	BtnSettings = MainMenu.Text(commands.CommandSettings)
	BtnBack     = MainMenu.Text(commands.CommandBack)

	// Room menu buttons.
	BtnShop     = RoomMenu.Text(commands.CommandShop)
	BtnAquaMan  = RoomMenu.Text(commands.CommandAquaMan)
	BtnCleanMan = RoomMenu.Text(commands.CommandCleanMan)

	// Aqua man menu buttons.
	BtnBringWater = AquaManMenu.Text(commands.CommandBringWater)

	// Clean man menu buttons.
	BtnCleanRoom = CleanManMenu.Text(commands.CommandCleanRoom)

	// Shop menu buttons.
	BtnCheckShopping  = ShopMenu.Text(commands.CommandCheck)
	BtnUploadPurchase = ShopMenu.Text(commands.CommandUploadPurchase)
	BtnPurchaseDone   = ShopMenu.Text(commands.CommandPurchaseDone)
	BtnPostPurchase   = ShopMenu.Text(commands.CommandPostPurchase)

	// News menu buttons.
	BtnNewsDone    = NewsMenu.Text(commands.CommandNewsDone)
	BtnPostNews    = NewsMenu.Text(commands.CommandPostNews)
	BtnDeleteDraft = NewsMenu.Text(commands.CommandDeleteDraft)

	// Setting menu buttons.
	BtnNotificationSettings = SettingsMenu.Text(commands.CommandNotificationSettings)
	BtnSettingsOfBot        = SettingsMenu.Text(commands.CommandSettingsOfBot)

	BtnUploadExam = ExamMenu.Text(commands.CommandUploadExam)
	BtnGetExam    = ExamMenu.Text(commands.CommandGetExam)
	BtnExamDone   = ExamMenu.Text(commands.CommandExamDone)
)

func InitializeMenus() {
	MainMenu.Reply(
		MainMenu.Row(BtnRoom, BtnNews),
		MainMenu.Row(BtnExam, BtnSettings),
	)

	RoomMenu.Reply(
		RoomMenu.Row(BtnShop, BtnAquaMan, BtnCleanMan),
		RoomMenu.Row(BtnBack),
	)

	AquaManMenu.Reply(
		AquaManMenu.Row(BtnBringWater),
		AquaManMenu.Row(BtnBack),
	)

	CleanManMenu.Reply(
		CleanManMenu.Row(BtnCleanRoom),
		CleanManMenu.Row(BtnBack),
	)

	ShopMenu.Reply(
		ShopMenu.Row(BtnUploadPurchase, BtnCheckShopping),
		ShopMenu.Row(BtnBack),
	)

	ExamMenu.Reply(
		ExamMenu.Row(BtnUploadExam, BtnGetExam),
		ExamMenu.Row(BtnBack),
	)

	NewsMenu.Reply(
		NewsMenu.Row(BtnNewsDone, BtnDeleteDraft),
		NewsMenu.Row(BtnBack),
	)

	PostNewsMenu.Reply(
		PostNewsMenu.Row(BtnPostNews, BtnDeleteDraft),
		PostNewsMenu.Row(BtnBack),
	)

	SettingsMenu.Reply(
		SettingsMenu.Row(BtnSettingsOfBot, BtnNotificationSettings),
		SettingsMenu.Row(BtnBack),
	)

	ShopUploadMenu.Reply(
		ShopUploadMenu.Row(BtnPurchaseDone, BtnDeleteDraft),
		ShopUploadMenu.Row(BtnBack),
	)

	PostPurchaseMenu.Reply(
		PostPurchaseMenu.Row(BtnPostPurchase, BtnDeleteDraft),
		PostPurchaseMenu.Row(BtnBack),
	)

	ExamUploadMenu.Reply(
		ExamUploadMenu.Row(BtnExamDone, BtnDeleteDraft),
		ExamUploadMenu.Row(BtnBack),
	)
}

func GetMenus() map[string]*telegram.ReplyMarkup {
	allMenus := map[string]*telegram.ReplyMarkup{}

	allMenus[commands.CommandStart] = MainMenu
	allMenus[commands.CommandRoom] = RoomMenu
	allMenus[commands.CommandExam] = ExamMenu
	allMenus[commands.CommandNews] = NewsMenu
	allMenus[commands.CommandSettings] = SettingsMenu
	allMenus[commands.CommandAquaMan] = AquaManMenu
	allMenus[commands.CommandCleanMan] = CleanManMenu
	allMenus[commands.CommandShop] = ShopMenu
	allMenus[commands.CommandNewsDone] = PostNewsMenu
	allMenus[commands.CommandUploadPurchase] = ShopUploadMenu
	allMenus[commands.CommandPurchaseDone] = PostPurchaseMenu
	allMenus[commands.CommandUploadExam] = ExamUploadMenu

	return allMenus
}
