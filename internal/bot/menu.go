package bot

import (
	"strings"

	tb "gopkg.in/telebot.v3"
)

const (
	next = "menu_next"
	back = "menu_back"
)

var (
	FirstMenu  = &tb.ReplyMarkup{}
	SecondMenu = &tb.ReplyMarkup{}

	btnNext = FirstMenu.Data("➡️ Далее", "menu_next")
	btnBack = SecondMenu.Data("⬅️ Назад", "menu_back")
	btnHelp = SecondMenu.URL("📘 Помощь", "Никто тебе не поможет =(")
)

// При старте — добавь row
func setupMenus() {
	// Первое меню: одна кнопка
	FirstMenu.Inline(
		FirstMenu.Row(btnNext),
	)

	// Второе меню: две кнопки в разных строках
	SecondMenu.Inline(
		SecondMenu.Row(btnBack),
		SecondMenu.Row(btnHelp),
	)
}

func (b *botHandle) showFirstMenu(c tb.Context) error {
	_, err := c.Bot().Send(c.Chat(), "<b>Меню 1</b>\n\nВыберите действие:", &tb.SendOptions{
		ParseMode:   tb.ModeHTML,
		ReplyMarkup: FirstMenu,
	})
	return err
}

func (b *botHandle) showSecondMenu(c tb.Context) error {
	_, err := c.Bot().Send(c.Chat(), "<b>Меню 2</b>\n\nБольше опций:", &tb.SendOptions{
		ParseMode:   tb.ModeHTML,
		ReplyMarkup: SecondMenu,
	})
	return err
}

func (b *botHandle) onCallback(c tb.Context) error {
	data := c.Callback().Data

	switch data {
	case next:
		return b.showSecondMenu(c)
	case back:
		return b.showFirstMenu(c)
	default:
		return c.Respond(&tb.CallbackResponse{
			Text: "Неизвестная команда",
		})
	}
}

func escapeMarkdownV2(text string) string {
	// Экранируем все специальные символы, кроме тех, что используются для разметки
	// Например, если * используется для жирного текста, не экранируем его
	// простой пример
	specialChars := []string{"_", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
