-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_type
(
    id                  serial,
    event_type_name     varchar(20)         not null,
    created_at          timestamp default   current_timestamp,
    updated_at          timestamp default   current_timestamp,
  	PRIMARY KEY(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_type;
-- +goose StatementEnd
