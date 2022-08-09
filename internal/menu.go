package internal

import telegram "gopkg.in/telebot.v3"

var (
	// Universal markup builders.
	menu = &telegram.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	btnRoom     = menu.Text("\U0001F6D6 Room")
	btnNews     = menu.Text("📰 News")
	btnExam     = menu.Text("🧐 Exam")
	btnSettings = menu.Text("⚙️ Settings")
)

func setupMenu() {
	menu.Reply(
		menu.Row(btnRoom),
		menu.Row(btnNews),
		menu.Row(btnExam),
		menu.Row(btnSettings),
	)
}
