package saluteSpeech

const (
	insertFile  = "https://smartspeech.sber.ru/rest/v1/data:upload"
	createTask  = "https://smartspeech.sber.ru/rest/v1/speech:async_recognize"
	checkStatus = "https://smartspeech.sber.ru/rest/v1/task:get"
	getResult   = "https://smartspeech.sber.ru/rest/v1/data:download"

	authURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
)

// Структура для парсинга ответа от /data:upload
type uploadResp struct {
	Result struct {
		RequestFileID string `json:"request_file_id"`
	} `json:"result"`
}

type taskResp struct {
	Result struct {
		ID string `json:"id"` // Это и есть task_id для дальнейших запросов
	} `json:"result"`
}

type TaskStatusResult struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}

type statusResp struct {
	Status int              `json:"status"`
	Result TaskStatusResult `json:"result"`
}
