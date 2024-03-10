-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS user_schema;

CREATE TABLE IF NOT EXISTS user_schema.users (
    user_id      BIGSERIAL PRIMARY KEY,
    tg_id BIGINT not null unique,
    last_mes varchar(400) not null,
    progress INTEGER not null,
    curr_quotation INTEGER not null
);

CREATE TABLE IF NOT EXISTS user_schema.users_info (
    info_id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    nick varchar(20),
    money BIGINT,
    FOREIGN KEY (user_id) REFERENCES user_schema.users(user_id)
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE user_schema.users_info;
DROP TABLE user_schema.users;
DROP SCHEMA user_schema;

-- +goose StatementEnd