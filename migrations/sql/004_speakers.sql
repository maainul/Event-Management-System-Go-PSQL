-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS speakers
(
    id         serial ,
    first_name varchar(20)        not null,
    last_name  varchar(20)        not null,
    phone      varchar(11)        not null,
    address    varchar(200)       not null,
    username   varchar(20)        unique not null,
    email      varchar(50)        unique not null,
    created_at timestamp default  current_timestamp,
    updated_at  timestamp default current_timestamp,
  
  PRIMARY KEY(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS speakers;
-- +goose StatementEnd
