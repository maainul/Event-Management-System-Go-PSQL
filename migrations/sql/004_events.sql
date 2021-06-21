-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(  
    id                       serial,
    event_name               VARCHAR(160),
    speakers_id              INT REFERENCES speakers(id),
    event_type_id            INT REFERENCES event_type(id),
    start_time         		   timestamp,
    end_time 			           timestamp,
    event_date 				       Date,
    number_of_guest          INT,
    per_person_price         INT,
    status                   boolean                   default true,
    created_at               timestamp default         current_timestamp,
    updated_at               timestamp default  current_timestamp,
  	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
