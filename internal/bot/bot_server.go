package bot

import (
	"TgNoteAssist/internal/service"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"time"

	tb "gopkg.in/telebot.v3"
)

type botHandle struct {
	TgBot     *service.TgAssist // внутриняя логика приложения
	BotServer *tb.Bot
	lg        *slog.Logger
	token     string
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewBot(TgBot *service.TgAssist, logger *slog.Logger, token string) *botHandle {
	return &botHandle{TgBot: TgBot, lg: logger, token: token}
}

// StartHTTP - инициализация и запуск сервера
func (b *botHandle) StartBot() (*tb.Bot, error) {
	// Инициализация бота с токеном и настройками поллинга
	bot, err := tb.NewBot(tb.Settings{
		Token:     b.token,
		Poller:    &tb.LongPoller{Timeout: 10 * time.Minute},
		ParseMode: tb.ModeHTML,
		//Client:    httpClient,
	})
	if err != nil {
		return nil, err
	}
	b.BotServer = bot
	// Регистрация обработчиков
	bot.Handle(tb.OnCallback, b.onCallback)
	bot.Handle(tb.OnText, b.onText)   // Обработчик текстовых сообщений (включая команды)
	bot.Handle(tb.OnVoice, b.onVoice) // Обработчик голосовых сообщений
	bot.Handle(tb.OnAudio, b.onAudio) // Обработчик аудиофайлов

	// Запуск бота
	go func() {
		b.lg.Info("StartBot.start - успешно запущен и готов к работе!")
		bot.Start()
	}()

	return bot, nil

}
