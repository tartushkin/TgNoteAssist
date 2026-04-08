package repository

import (
	"TgNoteAssist/internal/model"
	"context"
	"fmt"
)

func (r *Repo) RegisterUser(ctx context.Context, user *model.User) error {
	_, err := r.conn.ExecContext(ctx, `
	INSERT INTO tgassist.t_users(n_tg_id, s_first_name, s_user_name, dt_created_at)
	VALUES ($1,$2,$3,$4);
	`, user.TgID, user.FirstName, user.NickName, user.Created)
	if err != nil {
		return fmt.Errorf(" ошибка при вставке нового пользователя: %s", err.Error())
	}
	return nil
}

func (r *Repo) GetUserList(ctx context.Context) ([]*model.User, error) {
	rows, err := r.conn.QueryContext(ctx,
		`SELECT n_tg_id, s_first_name, s_user_name, dt_created_at
		 FROM tgassist.t_users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.TgID, &user.FirstName, &user.NickName, &user.Created); err != nil {
			return nil, err
		}
		list = append(list, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// GetOriginalURL - получить оригинальный url по алиасу.
func (r *Repo) GetTrnscMeet(ctx context.Context, nickName string, meetID int) (string, error) {
	var transcription string

	err := r.conn.QueryRowContext(ctx,
		`SELECT t_text
		FROM tgassist.t_meetings
		WHERE n_id = $1
		AND s_user = $2`,
		meetID, nickName).Scan(&transcription)
	if err != nil {
		return "", err
	}

	return transcription, nil
}
func (r *Repo) InsertQuestion(ctx context.Context, nickName, quest string) error {
	_, err := r.conn.ExecContext(ctx, `
	INSERT INTO tgassist.t_chat(s_ask, s_user)
	VALUES ($1,$2);
	`, quest, nickName)
	if err != nil {
		return fmt.Errorf(" ошибка при вставке нового пользователя: %s", err.Error())
	}
	return nil
}
