package repository

import (
	"TgNoteAssist/internal/model"
	"context"
	"fmt"
	"time"
)

// GetMeetsList -
func (r *Repo) GetFileList(ctx context.Context, nickName string) ([]*model.File, error) {
	rows, err := r.conn.QueryContext(ctx, `
	    SELECT s_file_name, dt_created_at, n_msg_id
		FROM tgassist.t_files
		WHERE s_user = $1
		ORDER BY dt_created_at DESC
		`, nickName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fileList := []*model.File{}
	for rows.Next() {
		file := &model.File{}
		if err := rows.Scan(&file.Name, &file.Created, &file.MsgID); err != nil {
			return nil, err
		}
		fileList = append(fileList, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return fileList, nil
}

// GeetFile -
func (r *Repo) GeetFile(ctx context.Context, msgID int) (string, error) {
	var text string
	row := r.conn.QueryRowContext(ctx, `
	    SELECT t_text
		FROM tgassist.t_files
		WHERE n_msg_id = $1
		`, msgID)

	err := row.Scan(&text)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (r *Repo) InsertFile(ctx context.Context, file *model.File, user string) error {
	_, err := r.conn.ExecContext(ctx, `
	INSERT INTO tgassist.t_files(
	n_msg_id,
    s_file_name  ,
    s_file_id    ,
    s_user       ,
    n_file_size  ,
    s_status     ,
    dt_created_at,
	t_text)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`, file.MsgID, file.Name, file.FileID, user, file.Size, file.Status, file.Created, file.Text)
	if err != nil {
		return fmt.Errorf(" ошибка при вставке нового файла: %s", err.Error())
	}

	return nil
}

func (r *Repo) UpdateFile(ctx context.Context, msgID int, status string) error {

	query := `
	    UPDATE tgassist.t_files
	    SET s_status =  $1, dt_update = $2
	    WHERE n_msg_id = $3
	`
	_, err := r.conn.ExecContext(ctx, query, status, time.Now(), msgID)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) InsertText(ctx context.Context, msgID int, text string) error {

	query := `
	    UPDATE tgassist.t_files
	    SET t_text =  $1, dt_update = $2
	    WHERE n_msg_id = $3
	`
	_, err := r.conn.ExecContext(ctx, query, text, time.Now(), msgID)
	if err != nil {
		return err
	}
	return nil
}

// В вашем слое repository (например, repo.go)
func (r *Repo) FindUserChats(ctx context.Context, nickName, keywords string) ([]*model.File, error) {
	// Используем ILIKE для регистронезависимого поиска
	query := `
        SELECT s_file_name, dt_created_at, n_msg_id
        FROM tgassist.t_files
        WHERE s_user = $1 AND t_text ILIKE '%' || $2 || '%'
        ORDER BY dt_created_at DESC
        LIMIT 10;`

	rows, err := r.conn.QueryContext(ctx, query, nickName, keywords)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	fileList := []*model.File{}
	for rows.Next() {
		file := &model.File{}
		if err := rows.Scan(&file.Name, &file.Created, &file.MsgID); err != nil {
			return nil, err
		}
		fileList = append(fileList, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return fileList, nil
}
