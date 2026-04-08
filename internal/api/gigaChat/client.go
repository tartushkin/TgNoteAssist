package gigaChat

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	//baseURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions" // генирация текста и изображения
	baseURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	authURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	//baseURL = "https://gigachat.devices.sberbank.ru/v1/chat/completions"
)

// GigaChatClient — клиент для работы с GigaChat API
type Client struct {
	Token   string
	client  *resty.Client
	authKey string

	lg *slog.Logger
}

// GigaChatResponse — структура для полного ответа от GigaChat API
type gigaChatResponse struct {
	Object  string `json:"object"`  // "chat.completion"
	Created int64  `json:"created"` // Метка времени Unix
	Model   string `json:"model"`   // Название модели, например "GigaChat"

	Choices []struct {
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"` // "stop", "length" и т.д.

		Message struct {
			Role    string `json:"role"`    // "assistant" или "user"
			Content string `json:"content"` // Текст ответа
		} `json:"message"`
	} `json:"choices"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	SystemTokens     int `json:"system_tokens"`
}

// NewGigaChatClient создаёт новый клиент
func NewGigaChatClient(ctx context.Context, lg *slog.Logger, pool *x509.CertPool, authKey string) (*Client, error) {

	gg := &Client{
		authKey: authKey,
		lg:      lg,
	}

	//pool, err := setCert()
	//if err != nil {
	//	return nil, err
	//}

	tlsConfig := &tls.Config{InsecureSkipVerify: false, RootCAs: pool}
	client := resty.New()
	client.SetTLSClientConfig(tlsConfig)
	gg.client = client

	//accessToken, err := gg.GetAccessToken()
	//if err != nil {
	//	return nil, err
	//}
	//gg.Token = accessToken

	//go gg.processGetToken(ctx, ping)
	return gg, nil
}

// 1. Получи access_token
func (gg *Client) GetAccessToken() error {

	var result struct {
		AccessToken string `json:"access_token"`
	}

	_, err := gg.client.
		R().
		SetHeader("Authorization", "Bearer "+gg.authKey).
		SetHeader("RqUID", uuid.New().String()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody("scope=GIGACHAT_API_PERS").
		SetResult(&result).
		Post(authURL)
	if err != nil {
		return err
	}
	gg.Token = result.AccessToken
	return nil
}

// Generate — отправляет запрос к модели
// Generate — отправляет запрос к модели
func (g *Client) Generate(prompt string) (string, error) {

	// --- ИЗМЕНЕНИЕ 1: Используем простую структуру для отладки ---
	// Мы не будем использовать сложную структуру GigaChatResponse,
	// а просто распарсим ответ в map, чтобы увидеть его целиком.
	var rawResponse map[string]interface{}

	jsonBody := fmt.Sprintf(`{
	    "model": "GigaChat-Pro-preview",
        "messages": [{"role": "user", "content": "%s"}],
        "temperature": 0.7,
        "max_tokens": 512,
        "stream": false
    }`, prompt)

	// --- ИЗМЕНЕНИЕ 2: Используем метод SetError ---
	// Этот метод позволяет поймать ответ сервера, даже если он содержит ошибку (например, 4xx или 5xx).
	// Но самое главное - он позволяет увидеть тело ответа в переменной err.
	_, err := g.client.R().
		SetHeader("Authorization", "Bearer "+g.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		SetResult(&rawResponse). // Парсим в map для гибкости
		SetError(&rawResponse).  // Если статус не 2xx, тоже парсим сюда
		Post(baseURL)

	// Теперь выведем ВСЁ, что у нас есть
	//fmt.Printf("--- ОТВЕТ СЕРВЕРА ---\n")
	//fmt.Printf("Статус код: %d\n", resp.StatusCode())
	//fmt.Printf("Тело ответа (сырое): %s\n", resp)
	//fmt.Printf("Тело ответа (распарсенное): %+v\n", rawResponse)
	//fmt.Println("--------------------")

	if err != nil {
		// Здесь мы можем проверить, нет ли в rawResponse ключа "error" или "detail"
		if errorMsg, ok := rawResponse["error"].(string); ok {
			return "", fmt.Errorf("ошибка от API: %s", errorMsg)
		}

		return "", fmt.Errorf("неизвестная ошибка при запросе: %v", err)
	}

	// Если мы здесь, значит статус был 2xx и rawResponse заполнен
	if len(rawResponse) == 0 {
		return "", fmt.Errorf("пустой ответ от GigaChat")
	}

	// Пытаемся достать ответ из распарсенного map
	choices, ok := rawResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("пустой или неверный ответ от GigaChat: %v", rawResponse)
	}

	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})

	answer, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("не удалось получить текст ответа")
	}

	return answer, nil
}
