
CREATE table IF NOT EXISTS tgassist.t_users (
    n_tg_id BIGINT NOT NULL,
    s_first_name VARCHAR(100) NOT NULL,
    s_user_name VARCHAR(200) NOT NULL PRIMARY KEY,
    dt_created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  tgassist.t_users IS 'Таблица пользователей';
COMMENT ON COLUMN tgassist.t_users.n_tg_id IS 'tg-id пользоваателя';
COMMENT ON COLUMN tgassist.t_users.s_first_name IS 'имя пользоваателя';
COMMENT ON COLUMN tgassist.t_users.s_user_name IS 'nickName пользоваателя';
COMMENT ON COLUMN tgassist.t_users.dt_created_at IS 'вермя регистрации пользоваателя в нашем боте ';
