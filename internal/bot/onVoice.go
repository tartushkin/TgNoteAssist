package bot

import (
	"TgNoteAssist/internal/model"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func (b *botHandle) onVoice(m tb.Context) error {

	b.lg.Info("onText.Start - началао обработки голосового от бота.")
	incomUser := m.Sender()
	user, ok := b.TgBot.GetUser(incomUser.ID)
	if !ok {
		b.lg.Info("onAudio.unknown - пользователь: " + incomUser.Username + " - неизвестен сервису.")
		return m.Reply(escapeMarkdownV2("*Доброго времени суток!*.\n Необходимо зарегистрироваться по команде *'/start'*- для начала работы с ботом."), tb.ModeMarkdownV2)
	}
	err := b.base(m, user, "voice")
	if err != nil {
		return m.Reply(escapeMarkdownV2(model.ErrMsg), tb.ModeMarkdownV2)
	}
	return m.Reply(escapeMarkdownV2("*Расшифровываем🔄*"), tb.ModeMarkdownV2)
}

func (b *botHandle) base(m tb.Context, user *model.User, fileType string) error {
	var tgFile tb.File
	var fileName string
	msg := m.Message()
	if fileType == "voice" {
		fileName = user.NickName + "_VOICE_" + strconv.Itoa(msg.ID)
		tgFile = msg.Voice.File
	} else {
		fileName = msg.Audio.FileName
		tgFile = msg.Audio.File
	}
	file, err := b.getFile(tgFile, fileName) // скачали файл
	if err != nil {
		b.lg.Error("onAudio.err - возникла ошибка: " + err.Error())
		return err
	}

	callback := func(chatID int64, msgID int, text string) error {
		if text == "" {
			text = model.ErrMsg
		}
		_, err := b.BotServer.Send(
			&tb.Chat{ID: chatID},
			escapeMarkdownV2(text),
			tb.Silent,
			&tb.SendOptions{
				ReplyTo: &tb.Message{ID: msgID},
			},
		)
		return err
	}

	err = b.TgBot.RegisterFile(file, user, tgFile, msg, callback, fileName)
	if err != nil {
		b.lg.Error("onAudio.err - возникла ошибка: " + err.Error())
		return m.Reply(escapeMarkdownV2(model.ErrMsg), tb.ModeMarkdownV2)
	}

	b.lg.Info("onAudio.Complete - успешная обработока аудио запроса от: " + user.NickName)
	return nil
}
