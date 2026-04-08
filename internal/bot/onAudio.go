package bot

import (
	"TgNoteAssist/internal/model"
	"fmt"
	"os"

	tb "gopkg.in/telebot.v3"
)

// получаем запрос
// забираем id файла --> создаем задачу/встречу --> ьросаем бфзу и кеш
// скачиваем файл(при успехе)--> запускаем асинхронный запрос на распознование --> отправляем ответ юзеру(распознаем)
// в рутине отправили файл --> запрашиваем статус--> после полуения ответа бросаем данные база/кеш --> отправляем msg юзеру
// ----
// todo:
// делать ли ретрай логику?
// процесс подхвата задач в незвершенном статусе?
func (b *botHandle) onAudio(m tb.Context) error {
	b.lg.Info("onAudio.Start - началао обработки аудио-файла от бота.")

	incomUser := m.Sender()
	user, ok := b.TgBot.GetUser(incomUser.ID)
	if !ok {
		b.lg.Info("onAudio.unknown - пользователь: " + incomUser.Username + " - неизвестен сервису.")
		return m.Reply(escapeMarkdownV2("*Доброго времени суток!*.\n Необходимо зарегистрироваться по команде *'/start'*- для начала работы с ботом."), tb.ModeMarkdownV2)
	}

	err := b.base(m, user, "audio")
	if err != nil {
		return m.Reply(escapeMarkdownV2(model.ErrMsg), tb.ModeMarkdownV2)
	}
	return m.Reply(escapeMarkdownV2("*Расшифровываем🔄*"), tb.ModeMarkdownV2)
}

func (b *botHandle) getFile(tgFile tb.File, fileName string) (*os.File, error) {

	out, err := os.CreateTemp("fileList", fileName+".oga")
	if err != nil {
		return nil, fmt.Errorf("GetFile.create - ошибка создания внутриннего файла: %w", err)
	}
	defer out.Close()

	err = b.BotServer.Download(&tgFile, out.Name())
	if err != nil {
		return nil, fmt.Errorf("GetFile.download - ошибка создания внутриннего файла: %w", err)
	}

	b.lg.Info("getFile.complete - Файл успешно загружен по пути: " + out.Name())
	return out, nil
}
