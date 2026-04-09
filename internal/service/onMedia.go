package service

import (
	"TgNoteAssist/internal/model"
	"os"
	"time"

	tb "gopkg.in/telebot.v3"
)

type CallbackFunc func(chatID int64, msgID int, text string) error

func (t *TgAssist) RegisterFile(file *os.File, user *model.User, tgFile tb.File, msg *tb.Message, response CallbackFunc, fileName string) error {
	t.Lg.Info("RegisterFile.start - формировование файла - " + fileName)

	f, err := t.createFile(user, msg, &tgFile, fileName)
	if err != nil {
		return err
	}
	t.Lg.Info("RegisterFile.done - файл успешно сформирован и упакован в базу/кеш: " + fileName)
	// запускаем рутину, которая в себе отправляет файл и дальше уходит на опрос статуса
	go t.getText(response, f, file)
	return nil
}

// createMeet - формирование файла и упаковка в базу/кеш
func (t *TgAssist) createFile(user *model.User, msg *tb.Message, tgFile *tb.File, fileName string) (*model.File, error) {

	file := model.File{
		Name:    fileName,
		FileID:  tgFile.FileID,
		MsgID:   msg.ID,
		Size:    tgFile.FileSize,
		Created: time.Now(),
		Status:  model.REGISTERED,
		ChatID:  msg.Chat.ID,
		Text:    "",
	}

	err := t.Repo.InsertFile(t.Ctx, &file, user.NickName)
	if err != nil {
		return nil, err
	}

	t.Mu.RLock()
	user.FileList[file.MsgID] = &file
	t.Mu.RUnlock()
	return &file, nil
}

func (t *TgAssist) getText(response CallbackFunc, f *model.File, file *os.File) {
	t.Lg.Info("getText.start - старт процесса  'извлечение транскрипции' по файлу : " + f.Name)
	var text string
	defer func() {
		response(f.ChatID, f.MsgID, text)
	}()

	reqFileID, err := t.salutClient.SetFile(file.Name())
	if err != nil {
		t.Lg.Error("getText.SetFile.err - возникла ошибка при отправки файла на извлечение транскрипции: " + err.Error())
		return
	}

	taskID, err := t.salutClient.CreateRecognitionTask(reqFileID)
	if err != nil {
		t.Lg.Error("getText.CreateRecognitionTask.err - возникла ошибка при постановки задачи: " + err.Error())
		return
	}
	f.TaskID = taskID

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		t.Lg.Info("getText.loop - ожидание новой итерации опроса статуса по транскрипции файла: " + f.Name)
		select {
		case <-t.Ctx.Done():
			t.Lg.Info("getText.ctxDone - контекст сервиса был завершен")
			return
		case <-ticker.C:
			ticker.Reset(t.paramGetStatus)

			status, id, err := t.salutClient.CheckStatus(taskID)
			if err != nil {
				text = "возникла ошибка в работе сервиса, попробуйдет позже"
				t.Lg.Error("getText.CheckStatus.err - возникла ошибка при опросе статуса задачи: " + err.Error())
				return
			}

			err = t.updateStatus(f, status)
			if err != nil {
				t.Lg.Error("updatefile.err - возникла ошибка при записи обновления файла: " + f.Name + " - в базе:" + err.Error())
				return
			}

			switch status {
			case model.RUNNING:
				t.Lg.Info("getText.process - задача в статусе: " + status + ", по файлу - " + f.Name)
				continue
			case model.DONE:
				t.Lg.Info("getText.completed - задача в статусе: " + status + ", по файлу - " + f.Name)
				transcript, err := t.salutClient.GetResult(id)
				if err != nil {
					t.Lg.Error("getText.GetResult.err - возникла ошибка при отправки файла на извлечение транскрипции: " + err.Error())
					return
				}
				t.Lg.Info("getText.completed - успешно получили расшифровку, по файлу - " + f.Name)
				text = transcript
				err = t.insertText(f, text)
				if err != nil {
					t.Lg.Error("updatefile.err - возникла ошибка при записи обновления файла: " + f.Name + " - в базе:" + err.Error())
					return
				}
				t.Lg.Info("getText.complete - успешно поулчили и обновили файл : " + f.Name)
				return
			case model.ERROR:
				t.Lg.Error("getText.err - прекращаем опрос,задача в статусе: " + status + ", по файлу - " + f.Name)
				return
			default:
				t.Lg.Error("getText.err - возникла ошибка при отправки файла на извлечение транскрипции: ")
				return
			}
		}
	}

}

func (t *TgAssist) updateStatus(file *model.File, status string) error {

	err := t.Repo.UpdateFile(t.Ctx, file.MsgID, status)
	if err != nil {
		return err
	}

	t.Mu.Lock()
	file.Status = status
	t.Mu.Unlock()

	// может кеш обновить
	t.Lg.Info("updateFile.comple - успешно обновили cтаус: " + status + " - у файла: " + file.Name)
	return nil
}

func (t *TgAssist) insertText(file *model.File, text string) error {
	err := t.Repo.InsertText(t.Ctx, file.MsgID, text)
	if err != nil {
		return err
	}

	t.Mu.Lock()
	file.Text = text
	t.Mu.Unlock()

	// может кеш обновить
	t.Lg.Info("updateFile.comple - успешно вставили транскрипцию файла: " + file.Name)
	return nil
}
