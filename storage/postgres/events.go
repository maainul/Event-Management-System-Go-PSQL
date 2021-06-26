package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
)

const e = `
SELECT 
	events.id, 
	event_name,
	event_type_name,
	start_time,
	end_time,
	event_date,
	number_of_guest,
	per_person_price,
	first_name,
	last_name
	   
FROM 
	events 

JOIN speakers ON 
	events.speakers_id = speakers.id

JOIN event_type ON 
	events.event_type_id = event_type.id;
`

func (s *Storage) GetEvent() ([]storage.Events, error) {
	event := make([]storage.Events, 0)
	if err := s.db.Select(&event, e); err != nil {
		return nil, err
	}

	fmt.Print(et)
	return event, nil
}

const createEventQuery = `
	INSERT INTO events(
		event_name,
		number_of_guest,
		per_person_price,
		event_date,
		start_time,
		end_time
	)
	VALUES(
		:event_name,
		:number_of_guest,
		:per_person_price,
		:event_date,
		:start_time,
		:end_time
	)
	RETURNING id
	`

func (s *Storage) CreateEvent(event storage.Events) (int32, error) {
	stmt, err := s.db.PrepareNamed(createEventQuery)
	if err != nil {
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, event); err != nil {
		return 0, err
	}
	return id, nil
}
