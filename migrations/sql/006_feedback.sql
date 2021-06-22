-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feedback
(  
    id                   serial,
    user_id              INT REFERENCES     users(id),
    message              VARCHAR(250),
    created_at           timestamp default  current_timestamp,
    updated_at           timestamp default  current_timestamp,
  	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedback;
-- +goose StatementEnd
