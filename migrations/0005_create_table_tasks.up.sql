CREATE table IF NOT EXISTS tgassist.t_files (
    n_msg_id         BIGINT PRIMARY KEY,
    s_file_name      VARCHAR(500),
    s_file_id        VARCHAR(400), --  непонятно нужен или нет, мы уже скачали его
    s_user           VARCHAR(400),
    n_file_size      BIGINT,
    s_status         VARCHAR(200),
    t_text           TEXT NOT NULL,
    dt_created_at   TIMESTAMPTZ DEFAULT NOW(),
    dt_update     TIMESTAMPTZ,
    FOREIGN KEY (s_user) REFERENCES tgassist.t_users(s_user_name) ON DELETE CASCADE
);

COMMENT ON TABLE  tgassist.t_files IS 'Таблица задач распознавания файла';
COMMENT ON COLUMN tgassist.t_files.n_msg_id IS 'идентификатор тг сообщения'; -- для 'в ответ'
COMMENT ON COLUMN tgassist.t_files.s_file_name IS 'имя файла'; -- для 'в ответ'
COMMENT ON COLUMN tgassist.t_files.s_file_id IS 'наименовани/id файла от тг';
COMMENT ON COLUMN tgassist.t_files.s_user IS 'ник пользователя';
COMMENT ON COLUMN tgassist.t_files.n_file_size IS 'размер файла в байтах';
COMMENT ON COLUMN tgassist.t_files.s_status IS 'статус распознавания';
COMMENT ON COLUMN tgassist.t_files.t_text IS 'статус распознавания';
COMMENT ON COLUMN tgassist.t_files.dt_created_at IS 'дата регистрации встречи в боте';
COMMENT ON COLUMN tgassist.t_files.dt_update IS 'дата успешного результата';


