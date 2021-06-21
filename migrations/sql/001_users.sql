-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users
(
    id         serial             not null,
    first_name varchar(20)        not null,
    last_name  varchar(20)        not null,
    username   varchar(20) unique not null,
    email      varchar(50) unique not null,
    password   varchar(255)       not null,
    is_active  boolean     default true,
    is_admin   boolean     default false,
    created_at timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp,

    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;