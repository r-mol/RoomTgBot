package menus

import (
	"RoomTgBot/internal/commands"

	telegram "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	MainMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	RoomMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ExamMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	NewsMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	SettingsMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	AquaManMenu      = &telegram.ReplyMarkup{ResizeKeyboard: true}
	InitAquaManMenu  = &telegram.ReplyMarkup{ResizeKeyboard: true}
	InitCleanManMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}
	CleanManMenu     = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ShopMenu         = &telegram.ReplyMarkup{ResizeKeyboard: true}

	PostNewsMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	ShopUploadMenu   = &telegram.ReplyMarkup{ResizeKeyboard: true}
	ShopCheckMenu    = &telegram.ReplyMarkup{ResizeKeyboard: true}
	PostPurchaseMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	ExamUploadMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}
	SubjectMenu    = &telegram.ReplyMarkup{ResizeKeyboard: true}

	ListMenu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	// Main menu buttons.
	BtnRoom     = MainMenu.Text(commands.CommandRoom)
	BtnNews     = MainMenu.Text(commands.CommandNews)
	BtnExam     = MainMenu.Text(commands.CommandExam)
	BtnSettings = MainMenu.Text(commands.CommandSettings)
	BtnBack     = MainMenu.Text(commands.CommandBack)
	BtnPrevious = ListMenu.Data(commands.CommandPrevious, "prev")
	BtnNext     = ListMenu.Data(commands.CommandNext, "next")
	BtnExit     = ListMenu.Data(commands.CommandExit, "exit")

	// Room menu buttons.
	BtnShop     = RoomMenu.Text(commands.CommandShop)
	BtnAquaMan  = RoomMenu.Text(commands.CommandAquaMan)
	BtnCleanMan = RoomMenu.Text(commands.CommandCleanMan)

	// Aqua man menu buttons.
	BtnBringWater   = AquaManMenu.Text(commands.CommandBringWater)
	BtnWaterIsOver  = AquaManMenu.Text(commands.CommandWaterIsOver)
	BtnBringWaterIN = InitAquaManMenu.Data(commands.CommandBringWater, "BW")
	BtnNotInInnoAQ  = InitAquaManMenu.Data(commands.CommandNotInInno, "NIIAQ")
	BtnCantAQ       = InitAquaManMenu.Data(commands.CommandCant, "CantAQ")

	// Clean man menu buttons.
	BtnCleanRoom   = CleanManMenu.Text(commands.CommandCleanRoom)
	BtnCleanRoomIN = InitCleanManMenu.Data(commands.CommandCleanRoom, "CR")
	BtnNotInInnoCR = InitCleanManMenu.Data(commands.CommandNotInInno, "NIICR")
	BtnCantCR      = InitCleanManMenu.Data(commands.CommandCant, "CantCR")

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

	// Exam menu buttons.
	BtnUploadExam = ExamMenu.Text(commands.CommandUploadExam)
	BtnGetExam    = ExamMenu.Text(commands.CommandGetExam)
	BtnExamDone   = ExamMenu.Text(commands.CommandExamDone)

	// Subjects menu buttons.
	Subject1  = SubjectMenu.Data("Compilers Construction", "CC")
	Subject2  = SubjectMenu.Data("Control Theory", "CT")
	Subject3  = SubjectMenu.Data("Data Mining", "DM")
	Subject4  = SubjectMenu.Data("Data Modeling and Databases I", "DMDI")
	Subject5  = SubjectMenu.Data("Data Modeling and Databases II", "DMDII")
	Subject6  = SubjectMenu.Data("Differential Equations", "DE")
	Subject7  = SubjectMenu.Data("Digital Signal Processing", "DSP")
	Subject8  = SubjectMenu.Data("Distributed Systems", "DS")
	Subject9  = SubjectMenu.Data("Fundamentals of Computer Security", "FCS")
	Subject10 = SubjectMenu.Data("Fundamentals of Robotics", "FR")
	Subject11 = SubjectMenu.Data("Game Theory", "GF")
	Subject12 = SubjectMenu.Data("Information Retrieval", "IR")
	Subject13 = SubjectMenu.Data("Information Theory", "IT")
	Subject14 = SubjectMenu.Data("Introduction to AI", "IA")
	Subject15 = SubjectMenu.Data("Introduction to Big Data", "IBD")
	Subject16 = SubjectMenu.Data("Introduction to Machine Learning", "IML")
	Subject17 = SubjectMenu.Data("Lean Software Development", "LSD")
	Subject18 = SubjectMenu.Data("Mechanics and Machines", "MM")
	Subject19 = SubjectMenu.Data("Networks", "N")
	Subject20 = SubjectMenu.Data("Network and Cyber Security", "NCS")
	Subject21 = SubjectMenu.Data("Non-Linear Optimization", "NLO")
	Subject22 = SubjectMenu.Data("Operating Systems", "OS")
	Subject23 = SubjectMenu.Data("Philosophy II", "PII")
	Subject24 = SubjectMenu.Data("Physics I", "PHYSICSI")
	Subject25 = SubjectMenu.Data("Practicum Project", "PP")
	Subject26 = SubjectMenu.Data("Probability and Statistics", "PS")
	Subject27 = SubjectMenu.Data("Robotics Systems", "RS")
	Subject28 = SubjectMenu.Data("Sensors and Sensing", "SS")
	Subject29 = SubjectMenu.Data("Software Architecture", "SA")
	Subject30 = SubjectMenu.Data("Software Project", "SP")
	Subject31 = SubjectMenu.Data("Software Systems Design", "SSD")
	Subject32 = SubjectMenu.Data("System and Network Administration", "SNA")
	Subject33 = SubjectMenu.Data("Theoretical Computer Science", "TCS")
	Subject34 = SubjectMenu.Data("Theoretical Mechanics", "TM")
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

	ExamMenu.Reply(
		ExamMenu.Row(BtnUploadExam, BtnGetExam),
		ExamMenu.Row(BtnBack),
	)

	NewsMenu.Reply(
		NewsMenu.Row(BtnNewsDone, BtnDeleteDraft),
		NewsMenu.Row(BtnBack),
	)

	SettingsMenu.Reply(
		SettingsMenu.Row(BtnSettingsOfBot, BtnNotificationSettings),
		SettingsMenu.Row(BtnBack),
	)

	AquaManMenu.Reply(
		AquaManMenu.Row(BtnBringWater, BtnWaterIsOver),
		AquaManMenu.Row(BtnBack),
	)

	InitAquaManMenu.Inline(
		InitAquaManMenu.Row(BtnNotInInnoAQ, BtnCantAQ),
		InitAquaManMenu.Row(BtnBringWaterIN),
	)

	InitCleanManMenu.Inline(
		InitCleanManMenu.Row(BtnNotInInnoCR, BtnCantCR),
		InitCleanManMenu.Row(BtnCleanRoomIN),
	)

	CleanManMenu.Reply(
		CleanManMenu.Row(BtnCleanRoom),
		CleanManMenu.Row(BtnBack),
	)

	ShopMenu.Reply(
		ShopMenu.Row(BtnUploadPurchase, BtnCheckShopping),
		ShopMenu.Row(BtnBack),
	)

	ShopUploadMenu.Reply(
		ShopUploadMenu.Row(BtnPurchaseDone, BtnDeleteDraft),
		ShopUploadMenu.Row(BtnBack),
	)

	ShopCheckMenu.Reply(
		ShopCheckMenu.Row(BtnBack),
	)

	PostPurchaseMenu.Reply(
		PostPurchaseMenu.Row(BtnPostPurchase, BtnDeleteDraft),
		PostPurchaseMenu.Row(BtnBack),
	)

	PostNewsMenu.Reply(
		PostNewsMenu.Row(BtnPostNews, BtnDeleteDraft),
		PostNewsMenu.Row(BtnBack),
	)

	ExamUploadMenu.Reply(
		ExamUploadMenu.Row(BtnExamDone, BtnDeleteDraft),
		ExamUploadMenu.Row(BtnBack),
	)

	ListMenu.Inline(ListMenu.Row(BtnPrevious, BtnExit, BtnNext))

	SubjectMenu.Inline(
		SubjectMenu.Row(Subject1),
		SubjectMenu.Row(Subject2),
		SubjectMenu.Row(Subject3),
		SubjectMenu.Row(Subject4),
		SubjectMenu.Row(Subject4),
		SubjectMenu.Row(Subject5),
		SubjectMenu.Row(Subject6),
		SubjectMenu.Row(Subject7),
		SubjectMenu.Row(Subject8),
		SubjectMenu.Row(Subject9),
		SubjectMenu.Row(Subject10),
		SubjectMenu.Row(Subject11),
		SubjectMenu.Row(Subject12),
		SubjectMenu.Row(Subject13),
		SubjectMenu.Row(Subject14),
		SubjectMenu.Row(Subject15),
		SubjectMenu.Row(Subject16),
		SubjectMenu.Row(Subject17),
		SubjectMenu.Row(Subject18),
		SubjectMenu.Row(Subject19),
		SubjectMenu.Row(Subject20),
		SubjectMenu.Row(Subject21),
		SubjectMenu.Row(Subject22),
		SubjectMenu.Row(Subject23),
		SubjectMenu.Row(Subject24),
		SubjectMenu.Row(Subject25),
		SubjectMenu.Row(Subject26),
		SubjectMenu.Row(Subject27),
		SubjectMenu.Row(Subject28),
		SubjectMenu.Row(Subject29),
		SubjectMenu.Row(Subject30),
		SubjectMenu.Row(Subject31),
		SubjectMenu.Row(Subject32),
		SubjectMenu.Row(Subject33),
		SubjectMenu.Row(Subject34),
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
	allMenus[commands.CommandUploadPurchase] = ShopUploadMenu
	allMenus[commands.CommandCheck] = ShopCheckMenu
	allMenus[commands.CommandNewsDone] = PostNewsMenu
	allMenus[commands.CommandPurchaseDone] = PostPurchaseMenu
	allMenus[commands.CommandUploadExam] = ExamUploadMenu
	allMenus[commands.CommandExamDone] = SubjectMenu
	allMenus[commands.CommandAquaManIN] = InitAquaManMenu
	allMenus[commands.CommandCleanManIN] = InitCleanManMenu

	return allMenus
}
