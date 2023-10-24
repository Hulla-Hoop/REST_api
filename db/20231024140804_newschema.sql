-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
id int NOT NULL,
created_at timestamp with time zone,
updated_at timestamp with time zone,
name text,
surname text,
patronymic text,
age integer,
gender text,
nationality text,
PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users ;
-- +goose StatementEnd
