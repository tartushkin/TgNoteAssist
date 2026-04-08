package bot

import (
	"TgNoteAssist/internal/model"
	"strings"

	tb "gopkg.in/telebot.v3"
)

const (
	start = "/start"
	get   = "/get"
	find  = "/find"
	chat  = "/chat"
	list  = "/list"
)

func (b *botHandle) onText(m tb.Context) error {
	// вытаскиваем/регаем пользователя
	// валидация текста(что за команда)
	// уходим на нужную логику

	b.lg.Info("onText.Start - началао обработки текста от бота.")
	incomUser := m.Sender()
	user, ok := b.TgBot.GetUser(incomUser.ID)
	if !ok {
		b.lg.Info("onAudio.unknown - пользователь: " + incomUser.Username + " - неизвестен сервису.")
		return m.Reply(escapeMarkdownV2("*Доброго времени суток!*.\n Необходимо зарегистрироваться по команде *'/start'*- для начала работы с ботом."), tb.ModeMarkdownV2)
	}
	user, res, err := b.validateCMD(m)
	if err != nil {
		b.lg.Error("onText.err - возникла ошибка: " + err.Error())
		m.Reply(escapeMarkdownV2(model.ErrMsg), tb.ModeMarkdownV2) //escapeMarkdownV2(
		return nil
	}

	var userNick string
	if user != nil {
		userNick = user.NickName
	} else {
		userNick = "неизвестный пользователь"
	}

	b.lg.Info("onText.Complete - успешная обработока текстового запроса от " + userNick)
	return m.Reply(escapeMarkdownV2(res), tb.ModeMarkdownV2)
}

func (b *botHandle) validateCMD(m tb.Context) (*model.User, string, error) {
	text := m.Text()
	incomUser := m.Sender()
	// Проверяем, является ли сообщение командой
	if strings.HasPrefix(text, "/") {
		// Разбиваем команду и аргументы (например, "/get 123")
		parts := strings.SplitN(text, " ", 2)
		command := parts[0]
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}
		if command == start {
			user, err := b.TgBot.RegisterUser(incomUser)
			if err != nil {
				return nil, "", err
			}
			return user, "Здравствуй 👋🎉, " + user.NickName + "!", nil
		}
		user, ok := b.TgBot.GetUser(incomUser.ID)
		if !ok {
			b.lg.Info("onAudio.unknown - пользователь: " + incomUser.Username + " - неизвестен сервису.")
			return nil, "*Доброго времени суток!*.\n Необходимо зарегистрироваться по команде *'/start'*- для начала работы с ботом.", nil
		}
		// Обрабатываем команды
		switch command {
		case list: //  список встреч
			result, err := b.TgBot.GetList(user)
			if err != nil {
				return user, "", err
			}
			return user, result, nil

		case get:
			// Логика получения текста встречи по id
			result, err := b.TgBot.Get(args)
			if err != nil {
				return nil, "", err
			}
			return user, result, nil

		case find:
			err := b.TgBot.Find(user.TgID, args)
			if err != nil {
				return user, "", err
			}
			return user, "result", nil

		case chat:
			res, err := b.TgBot.AskGigaChat(user, args)
			if err != nil {
				return user, "", err
			}
			return user, res, nil
		default: // res - string
			return nil, "*Неизвестная команда🚨👎😢\nИспользуйте /start, /list, /get, /find или  /chat*", nil
		}
	}
	result := "*Неизвестная команда*🚨👎😢\n*Используйте* /start, /list, /get, /find *или*  /chat"
	return nil, result, nil
}

//👋 (\U0001F44B) — "Привет"
//🙌 (\U0001F64C) — "Спасибо"
//👍 (\U0001F44D) — "Лайк"
//👎 (\U0001F44E) — "Дизлайк"
//
//
//Эмоции:
//
//😊 (\U0001F60A) — Улыбка
//😢 (\U0001F622) — Грусть
//😍 (\U0001F60D) — Влюблённость
//😱 (\U0001F631) — Удивление
//😠 (\U0001F620) — Злость
//
//
//2. Эмодзи для уведомлений
//
//
//Успех:
//
//✅ (\u2705) — Готово
//🎉 (\U0001F389) — Празднование
//🔥 (\U0001F525) — Горячо
//
//
//Ошибки:
//
//❌ (\u274C) — Ошибка
//⚠️ (\u26A0) — Предупреждение
//🚨 (\U0001F6A8) — Сигнал тревоги
//
//
//3. Эмодзи для команд
//
//
//Списки и меню:
//
//📋 (\U0001F4CB) — Список
//📝 (\U0001F4DD) — Заметка
//🔍 (\U0001F50D) — Поиск
//
//
//Действия:
//
//🔄 (\U0001F504) — Обновление
//⏳ (\u23F3) — Ожидание
//🗑️ (\U0001F5D1) — Удаление
//
//
//4. Эмодзи для времени и дат
//
//Время:
//
//⏰ (\u23F0) — Будильник
//🕒 (\U0001F552) — Часы
//📅 (\U0001F4C5) — Календарь
//
//
//5. Эмодзи для файлов и медиа
//
//Файлы:
//
//📄 (\U0001F4C4) — Документ
//🖼️ (\U0001F5BC) — Картинка
//🎵 (\U0001F3B5) — Музыка
//
//
//6. Эмодзи для общения
//
//Чат:
//
//💬 (\U0001F4AC) — Реплика
//🗣️ (\U0001F5E3) — Разговор
//📩 (\U0001F4E9) — Письмо
//
//
//7. Эмодзи для статусов
//
//Статусы:
//
//🆕 (\U0001F195) — Новое
//🔄 (\U0001F504) — Обновлено
//🔒 (\U0001F512) — Закрыто
