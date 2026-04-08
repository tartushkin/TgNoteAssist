CREATE table IF NOT EXISTS tgassist.t_chat (
    n_id     serial PRIMARY KEY,
    s_user   VARCHAR(100) NOT NULL,
    s_ask    text NOT NULL,
    dt_created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (s_user) REFERENCES tgassist.t_users(s_user_name) ON DELETE CASCADE
);

COMMENT ON TABLE  tgassist.t_chat IS 'Таблица встреч';
COMMENT ON COLUMN tgassist.t_chat.n_id IS 'идентификатор вопроса';
COMMENT ON COLUMN tgassist.t_chat.s_user IS 'пользователь';
COMMENT ON COLUMN tgassist.t_chat.s_ask IS 'вопрос к чату';
COMMENT ON COLUMN tgassist.t_chat.dt_created_at IS 'дата регистрации встречи в боте';
