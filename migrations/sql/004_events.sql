-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (   
  "id"                       serial,   
  "event_name"               VARCHAR(160),
  "event_type_id"            INT DEFAULT NULL,     
  "speakers_id"              INT DEFAULT NULL,
  "event_start_time"         time,
  "event_end_time" 			     time,
  "event_date" 				       date,
  "number_of_guest"          INT,
  "per_person_price"         INT,
  "status"                   boolean                   default true,
  "created_at"               timestamp default         current_timestamp,
  "updated_at"               timestamp default         current_timestamp,
  PRIMARY KEY ("id"),   
  FOREIGN KEY ("event_type_id") REFERENCES event_type("id") on delete cascade on update cascade,   
  FOREIGN KEY ("speakers_id") REFERENCES speakers("id") on delete cascade on update cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd

