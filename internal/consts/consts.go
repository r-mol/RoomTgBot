package consts

const (
	// Commands
	CommandStart       = "/start"
	CommandContact     = "/contact"
	CommandHelp        = "/help"
	CommandBringWater  = "/bring_water"
	CommandWaterIsOver = "/water_is_over"
	CommandCleanRoom   = "/clean_room"

	// Main consts
	CommandRoom     = "\U0001F6D6 Room"
	CommandNews     = "📰 News"
	CommandExam     = "🧐 Exam"
	CommandSettings = "⚙️ Settings"
	CommandBack     = "Back ↩️"

	// Room consts
	CommandShop       = "🛍 Shop"
	CommandAquaMan    = "\U0001F9A6 Aqua-Man"
	CommandCleanMan   = "\U0001F9F9 Clean-Man"
	CommandAquaManIN  = "Bring Water"
	CommandCleanManIN = "Clean Room"
	CommandNotInInno  = "🛫 Not in Innopolis"
	CommandCant       = "😥 Can't do it now"

	// Upload purchase consts
	CommandPostPurchase   = "📦 Post purchase"
	CommandUploadPurchase = "📥 Upload purchase"
	CommandPurchaseDone   = "✅ Purchase done"
	CommandCheckPurchases = "\U0001F9FE Check purchases"

	// News consts
	CommandPostNews    = "✉️ Post News"
	CommandUploadNews  = "📥 Upload News"
	CommandNewsDone    = "✅ News done"
	CommandCheckNews   = "\U0001F9FE Check News"
	CommandDeleteDraft = "🗑 Remove draft"

	// Settings consts
	CommandNotificationSettings = "📬 Notification settings"
	CommandSettingsOfBot        = "🔐 Settings of bot"

	// Exam consts
	CommandUploadExam = "📥 Upload exam"
	CommandGetExam    = "📤 Get exam"
	CommandExamDone   = "✅ Exam done"

	// Notification const
	Notification         = "Notification"
	NotificationNews     = "News"
	NotificationMoney    = "Money"
	NotificationExam     = "Exam"
	NotificationShop     = "Shop"
	NotificationCleaning = "Cleaning"

	InitState           = "init_state"
	BaseForConvertToInt = 10
	TimeOutMultiplier   = 10

	// DB consts
	MongoDBName = "RoomTgBot"
)
