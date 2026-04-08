package model

import (
	"encoding/json"
	"time"
)

const (
	CONFLICT    = "данный URL"
	ERRCONFLICT = "ERRCONFLICT"
	//  статус расчёта начисления из blackBox
	REGISTERED = "REGISTERED" // для новых
	RUNNING    = "RUNNING"
	NEW        = "NEW"
	ERROR      = "ERROR"
	DONE       = "DONE"
	CANCELED   = "CANCELED"

	ErrMsg = "*УПС 😱\nОшибка в сервисе😢, попробуйте позже.*"
)

type User struct {
	TgID      int64
	FirstName string
	NickName  string
	Created   string
	FileList  map[int]*File //[]*Order
}
type File struct {
	Name    string
	MsgID   int
	FileID  string
	TaskID  string
	Size    int64
	Status  string
	Created time.Time
	Text    string
	ChatID  int64
}
type AccrualRes struct {
	Order   string      `json:"order"`
	Status  string      `json:"status"`
	Accrual json.Number `json:"accrual"`
}
