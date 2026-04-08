package saluteSpeech

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Client struct {
	Token   string
	client  *resty.Client
	authKey string

	lg *slog.Logger
}

// NewGigaChatClient создаёт новый клиент
func NewSaluteClient(ctx context.Context, lg *slog.Logger, pool *x509.CertPool, authKey string) (*Client, error) {

	sc := &Client{
		authKey: authKey,
		lg:      lg,
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: false, RootCAs: pool}
	client := resty.New()
	client.SetTLSClientConfig(tlsConfig)
	sc.client = client

	//accessToken, err := gg.GetAccessToken()
	//if err != nil {
	//	return nil, err
	//}
	//gg.Token = accessToken

	//go gg.processGetToken(ctx, ping)
	return sc, nil
}

// // 1. Получи access_token
// func (gg *Client) GetAccessToken() (string, error) {
//
//		var result struct {
//			AccessToken string `json:"access_token"`
//		}
//
//		_, err := gg.client.
//			R().
//			SetHeader("Authorization", "Bearer "+gg.authKey).
//			SetHeader("RqUID", uuid.New().String()).
//			SetHeader("Content-Type", "application/x-www-form-urlencoded").
//			SetBody("scope=GIGACHAT_API_PERS").
//			SetResult(&result).
//			Post(authURL)
//		if err != nil {
//			return "", err
//		}
//
//		return result.AccessToken, nil
//	}
//
// getSaluteSpeechAccessToken получает access_token для API SaluteSpeech.
// clientID и clientSecret - это ваши учетные данные (логин:пароль) для сервиса.
func (sp *Client) GetAccessToken() error {
	var result struct {
		AccessToken string `json:"access_token"`
	}

	// 4. Отправляем POST-запрос на эндпоинт авторизации
	resp, err := sp.client.R().
		SetHeader("Authorization", "Basic "+sp.authKey). // Используем Basic Auth
		SetHeader("RqUID", uuid.New().String()).         // Заголовок из документации
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody("scope=SALUTE_SPEECH_PERS").
		SetResult(&result). // Указываем структуру для парсинга JSON-ответа
		Post(authURL)

	if err != nil {
		return fmt.Errorf("сетевая ошибка при получении токена: %v", err)
	}

	// 5. Проверяем статус ответа
	if resp.IsError() {
		// Если API вернуло ошибку (например, неверный логин/пароль), выводим детали
		return fmt.Errorf("ошибка авторизации: %s. Статус: %d", resp.Status(), resp.StatusCode())
	}

	// 6. Проверяем, что токен действительно пришел
	if result.AccessToken == "" {
		return fmt.Errorf("получен пустой access_token")
	}
	sp.Token = result.AccessToken
	return nil
}

// SendFileForRecognition отправляет файл на сервер и возвращает идентификатор задачи.
func (c *Client) SetFile(filePath string) (string, error) {
	// UploadFile загружает файл на сервер и возвращает request_file_id.
	// Этот ID затем используется для создания задачи распознавания.
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении файла в буфер: %v", err)
	}

	res := uploadResp{}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.Token).
		SetHeader("Content-Type", "audio/mpeg").
		SetHeader("Accept", "application/json").
		SetBody(fileBytes).
		SetResult(&res).
		Post(insertFile)

	if err != nil {
		return "", fmt.Errorf("ошибка при отправке файла: %v", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("API вернуло ошибку %d при загрузке файла: %s", resp.StatusCode(), resp.Status())
	}

	if res.Result.RequestFileID == "" {
		return "", fmt.Errorf("не удалось получить request_file_id из ответа")
	}
	return res.Result.RequestFileID, nil
}

// CreateRecognitionTask создает задачу на распознавание, используя ранее загруженный файл.
func (c *Client) CreateRecognitionTask(requestFileID string) (string, error) {
	requestBody := map[string]interface{}{
		"request_file_id": requestFileID,
		"options": map[string]interface{}{
			"model":          "general",
			"audio_encoding": "MP3",
		},
	}

	res := taskResp{}

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&res).
		Post(createTask)

	if err != nil {
		return "", fmt.Errorf("ошибка при создании задачи: %v", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("API вернуло ошибку %d при создании задачи: %s - %s", resp.StatusCode(), resp.Status(), string(resp.Body()))
	}

	if res.Result.ID == "" {
		return "", fmt.Errorf("не удалось получить id задачи из ответа")
	}

	return res.Result.ID, nil
}

// CheckTaskStatus проверяет статус задачи распознавания.
func (c *Client) CheckStatus(taskID string) (string, error) {
	// URL для проверки статуса задачи
	// 1. Формируем базовый URL

	urlWithQuery := fmt.Sprintf("%s?id=%s", checkStatus, taskID)

	res := statusResp{}
	// Отправка запроса
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.Token).
		SetResult(&res).
		Get(urlWithQuery)

	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса: %v", err)
	}
	//// Читаем "сырое" тело ответа
	//defer resp.RawBody().Close()
	//bodyBytes, _ := io.ReadAll(resp.RawBody())
	//rawResponse = string(bodyBytes)
	//
	//// Выводим в лог всё, что получили
	//fmt.Printf("--- ОТВЕТ СЕРВЕРА (СТАТУС) ---\n")
	//fmt.Printf("Код HTTP: %d\n", resp.StatusCode())
	//fmt.Printf("Тело ответа: %s\n", rawResponse)
	//fmt.Println("----------------------------------")
	//fmt.Println("ответ - ", res)
	if resp.IsError() {
		return "", fmt.Errorf("API вернуло ошибку %d: %s", resp.StatusCode(), resp.Status())
	}

	return res.Result.Status, nil
}

// GetTranscriptionResult получает финальный результат распознавания.
func (c *Client) GetResult(responseFileID string) (string, error) {
	// Формируем URL с query-параметром
	url := fmt.Sprintf("https://smartspeech.sber.ru/rest/v1/data:download?response_file_id=%s", responseFileID)

	var resultText string

	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.Token).
		SetHeader("Accept", "application/octet-stream").
		SetDoNotParseResponse(true).
		Get(url)

	if err != nil {
		return "", fmt.Errorf("ошибка при отправке запроса на скачивание: %v", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("API вернуло ошибку %d при скачивании: %s", resp.StatusCode(), resp.Status())
	}

	// Читаем "сырое" тело ответа (бинарные данные)
	defer resp.RawBody().Close()
	bytes, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении тела ответа: %v", err)
	}

	resultText = string(bytes)

	if resultText == "" {
		return "", fmt.Errorf("пустой ответ от сервера: файл с транскрипцией может быть пустым")
	}

	return resultText, nil
}
