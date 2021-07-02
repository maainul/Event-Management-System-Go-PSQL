-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS booking(
    "id"                    serial,
    "event_id" 			        INT DEFAULT NULL,
    "user_id"               INT DEFAULT NULL,
    "number_of_ticket"      INT,
    "total_amount"          INT,
    "created_at"            timestamp default         current_timestamp,
    "updated_at"            timestamp default         current_timestamp,
  PRIMARY KEY(id),
  CONSTRAINT event_id FOREIGN KEY(event_id) REFERENCES events(id),  
  CONSTRAINT user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booking;
-- +goose StatementEnd
