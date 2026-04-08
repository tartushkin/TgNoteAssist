package service

import (
	"TgNoteAssist/internal/model"
	"fmt"
	"strconv"
	"strings"
)

// GetList - список сохраненных встреч.
func (t *TgAssist) GetList(user *model.User) (string, error) {

	var result strings.Builder
	if len(user.FileList) < 1 {
		result.WriteString("*У вас пока еще нет встреч*😢")
		return result.String(), nil
	}

	result.WriteString("📋 *Список ваших встреч:*\n\n")

	for i, file := range user.FileList {
		result.WriteString(fmt.Sprintf(
			"%d. **%s**\n   🗓 *Дата:* %s\n   📝 *ID-встречи:* %v\n\n",
			i+1,
			file.Name,
			file.Created,
			file.MsgID,
		))
	}

	return result.String(), nil
}

// Get – получение текста встречи по id.
func (t *TgAssist) Get(msgID string) (string, error) {

	id, err := strconv.Atoi(msgID)
	if err != nil {
		return "", err
	}

	text, err := t.Repo.GeetFile(t.Ctx, id)
	if err != nil {
		return "", err
	}

	return text, nil
}

// Find - поиск встречи по ключевым словам.
func (t *TgAssist) Find(tgUserID int64, keywords string) error {

	// продумать логику
	return nil
}

// AskGigaChat - запрос к GigaChat.
func (t *TgAssist) AskGigaChat(user *model.User, request string) (string, error) {

	err := t.Repo.InsertQuestion(t.Ctx, user.NickName, request)
	if err != nil {
		return "", fmt.Errorf("askGigaChat - ошибка при записи в базу: %w", err)
	}
	res, err := t.gigaClient.Generate(request)
	if err != nil {
		return "", err
	}

	// продумать логику
	return res, nil
}
