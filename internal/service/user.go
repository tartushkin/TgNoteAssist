package service

import (
	"TgNoteAssist/internal/model"
	"context"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (m *TgAssist) LoadStorageUser(ctx context.Context) error {
	list, err := m.Repo.GetUserList(ctx)
	if err != nil {
		return err
	}

	for _, user := range list {
		fileList, err := m.Repo.GetFileList(ctx, user.NickName)
		if err != nil {
			return err
		}

		user.FileList = make(map[int]*model.File)
		for _, file := range fileList {
			user.FileList[file.MsgID] = file
		}
		m.Mu.Lock()
		m.UserCH[user.TgID] = user
		m.Mu.Unlock()
	}
	return nil
}

func (t *TgAssist) RegisterUser(inсomUser *tb.User) (*model.User, error) {
	user, ok := t.GetUser(inсomUser.ID)
	if ok {
		t.Lg.Info("RegisterUser - пользователь: " + inсomUser.Username + " - известен сервису.")
		return user, nil
	}

	t.Lg.Info("RegisterUser - пользователь: " + inсomUser.Username + " - неизвестен сервису. Старт регистрации пользователя.")
	user, err := t.registerUser(inсomUser)
	if err != nil {
		return nil, err
	}
	t.Lg.Info("RegisterUser.complete - успешно зарегистрировали новго пользователя в сервисе.")
	return user, nil
}

// chekUser - проверка наличия пользователя в сиситеме.
func (t *TgAssist) GetUser(incomID int64) (*model.User, bool) {
	t.Mu.RLock()
	user, ok := t.UserCH[incomID]
	t.Mu.RUnlock()
	if ok {
		return user, true
	}
	return nil, false
}

// registerUser - регистрация нового пользователя.
func (t *TgAssist) registerUser(inсomUser *tb.User) (*model.User, error) {
	tc := time.Now().Format(time.RFC3339)
	u := &model.User{
		TgID:      inсomUser.ID,
		FirstName: inсomUser.FirstName,
		NickName:  inсomUser.Username,
		Created:   tc,
	}

	err := t.Repo.RegisterUser(t.Ctx, u)
	if err != nil {
		return nil, err
	}

	t.Mu.Lock()
	t.UserCH[u.TgID] = u
	t.Mu.Unlock()

	return u, nil
}
