package service

import (
	"TgNoteAssist/internal/api/gigaChat"
	"TgNoteAssist/internal/api/saluteSpeech"
	cfg "TgNoteAssist/internal/config/app"
	"TgNoteAssist/internal/config/db"
	"TgNoteAssist/internal/model"
	"TgNoteAssist/internal/repository"
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"
)

type TgAssist struct {
	Lg          *slog.Logger
	Ctx         context.Context
	conn        *sql.DB
	Repo        *repository.Repo
	gigaClient  *gigaChat.Client
	salutClient *saluteSpeech.Client

	HTTPPort    string
	PathStorage string
	DNS         string
	UserCH      map[int64]*model.User

	Mu sync.RWMutex

	paramGetActualGigaToken time.Duration
	paramGetStatus          time.Duration
}

// Create - заполнение структуры приложения ,
func Create(ctx context.Context, Lg *slog.Logger, cfg *cfg.Config, ggcl *gigaChat.Client, spcl *saluteSpeech.Client) (*TgAssist, error) {
	userCH := map[int64]*model.User{}
	tg := &TgAssist{
		Lg:                      Lg,
		Ctx:                     ctx,
		UserCH:                  userCH,
		HTTPPort:                cfg.Port,
		gigaClient:              ggcl,
		salutClient:             spcl,
		paramGetActualGigaToken: cfg.Param_Get_Giga_Token,
		paramGetStatus:          cfg.Param_get_status,
	}

	if cfg.DbDNS != "" {
		conn, err := db.NewConnection(ctx, cfg.DbDNS)
		if err != nil {
			Lg.Error("db: не удалось подключиться к DB")
			return nil, err
		}
		tg.DNS = cfg.DbDNS
		Lg.Info("db: успешно подключились к DB")
		tg.conn = conn
		tg.Repo = repository.NewRepository(tg.conn)
	}

	err := tg.init(ctx, Lg)
	if err != nil {
		return nil, err
	}
	//m.checkUnfinishedOrder(ctx)
	go tg.processGetToken(ctx, cfg.Param_Get_Giga_Token)
	return tg, nil
}

// init - инициализация данных
func (tg *TgAssist) init(ctx context.Context, Lg *slog.Logger) error {
	Lg.Info("init.start - инициализация пользователей ")
	err := tg.LoadStorageUser(ctx)
	if err != nil {
		Lg.Error("init.error - ошибка загрузки пользователей: " + err.Error())
		return err
	}
	Lg.Info("init.finish - пользователи успешно подгружены")
	return nil
}

// processGetStatus - процес получения статуса и баллов по заказу
func (tg *TgAssist) processGetToken(ctx context.Context, param time.Duration) {
	tg.Lg.Info("processGetToken.start - старт процесса получения актуального токена")
	ticker := time.NewTicker(time.Second * 1) // запускаем процесс по тику

	for {
		tg.Lg.Info("processGetToken.wait - ожидание новой итерации  получения актуального токена " + param.String())
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(param * time.Minute)
			var err error
			for i := 0; i < 3; i++ {
				err = tg.gigaClient.GetAccessToken()
				if err == nil {
					break
				}
				tg.Lg.Warn("processGetToken.err - попытка получить токен провалилась, повтор через 1 сек...", "попытка", i+1)
				time.Sleep(time.Second)
			}
			for i := 0; i < 3; i++ {
				err = tg.salutClient.GetAccessToken()
				if err == nil {
					break
				}
				tg.Lg.Warn("processGetToken.err - попытка получить токен провалилась, повтор через 1 сек...", "попытка", i+1)
				time.Sleep(time.Second)
			}
			if err != nil {
				tg.Lg.Error("processGetToken.err - не удалось получить токен после 3 попыток: " + err.Error() + ". Уходим на некст итерацию по периоду.")
				continue
			}
			tg.Lg.Info("processGetToken.complete - токен получен.")
		}
	}
}
