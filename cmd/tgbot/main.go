package main

import (
	"TgNoteAssist/internal/api/gigaChat"
	"TgNoteAssist/internal/api/saluteSpeech"
	"TgNoteAssist/internal/bot"
	cfg "TgNoteAssist/internal/config/app"
	srv "TgNoteAssist/internal/service"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//ctx := context.Background()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := cfg.NewConfig() // инициализация конфига
	if err != nil {
		lg.Error("Ошибка инициализации GigaChat-клиента", "error", err)
		panic(err)
	}
	lg.Info("Конфиг загружен", "port", cfg.Port)

	ggcl, err := gigaChat.NewGigaChatClient(ctx, lg, cfg.PoolCert, cfg.GigaAuthKey)
	if err != nil {
		lg.Error("Ошибка инициализации GigaChat-клиента", "error", err)
		panic(err)
	}
	lg.Info("Клиент к gigaChat создан")

	spcl, err := saluteSpeech.NewSaluteClient(ctx, lg, cfg.PoolCert, cfg.GigaAuthKey)
	if err != nil {
		lg.Error("Ошибка инициализации SaluteSpeech-клиента", "error", err)
		panic(err)
	}
	lg.Info("Клиент к saluteSpeech создан")

	tgb, err := srv.Create(ctx, lg, cfg, ggcl, spcl) // инициализация сервиса
	if err != nil {
		lg.Error("Ошибка инициализации сервиса", "error", err)
		panic(err)
	}
	lg.Info("Сервис создан")

	h := bot.NewBot(tgb, tgb.Lg, cfg.BotToken) // бросать токен вместо порта
	lg.Info("Хендлеры созданы")
	lg.Info("Запуск сервера")
	go func() {
		if _, err := h.StartBot(); err != nil && err != http.ErrServerClosed {
			lg.Error("ошибка TG-сервера", "error", err)
			stop() // ускоряем остановку
		}
	}()

	lg.Info("HTTP-сервер запущен", "port", cfg.Port)

	// Ждём сигнала остановки
	<-ctx.Done()
	lg.Info("получен сигнал остановки, завершаем работу...")

}
