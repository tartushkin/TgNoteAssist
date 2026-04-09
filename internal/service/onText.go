package service

import (
	"TgNoteAssist/internal/model"
	"fmt"
	"iter"
	"strconv"
	"strings"
)

// GetList - список сохраненных встреч.
func (t *TgAssist) GetList(user *model.User) (string, error) {
	var lines []string

	if len(user.FileList) < 1 {
		return "*У вас пока еще нет встреч*😢", nil
	}

	lines = append(lines, "📋 *Список ваших встреч:*\n\n")

	for line := range getListItem(user.FileList) {
		lines = append(lines, line)
	}
	result := strings.Join(lines, "\n")
	return result, nil
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
func (t *TgAssist) Find(user, keywords string) (string, error) {
	// 1. Получаем данные из базы
	list, err := t.Repo.FindUserChats(t.Ctx, user, keywords)
	if err != nil {
		return "", fmt.Errorf("ошибка при поиске: %w", err)
	}

	if len(list) == 0 {
		return fmt.Sprintf("*Не нашли подходящих встреч по запросу:* *%s*\n*Попробуйте изменить ключевые слова.*", keywords), nil
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("📋 *Нашли совпадения по запросу '%s':*", keywords))

	for line := range generateChatLines(list) {
		lines = append(lines, line)
	}

	fullResult := strings.Join(lines, "\n")

	return fullResult, nil
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

// getListItem - это функция-итератор (генератор строк).
func getListItem(fileList map[int]*model.File) iter.Seq[string] {
	return func(yield func(string) bool) {
		var n int
		for i, file := range fileList {
			line := fmt.Sprintf(
				"%d. **%s**\n   🗓 *Дата:* %s\n   📝 *ID-встречи:* %v\n\n",
				n+1,
				file.Name,
				file.Created,
				i,
			)
			n++
			if !yield(line) {
				return
			}
		}
	}
}

// generateChatLines - генератор строк для списка встреч.
func generateChatLines(list []*model.File) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i, file := range list {
			line := fmt.Sprintf(
				"%d. **%s**\n   🗓 *Дата:* %s\n   📝 *ID-встречи:* %v\n\n",
				i+1,
				file.Name,
				file.Created.Format("02.01.2006"),
				file.MsgID,
			)
			if !yield(line) {
				return
			}
		}
	}
}
