package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"log"
	"strconv"
)

const e = `
SELECT 
	events.id, 
	event_name,
	event_type_name,
	event_start_time,
	event_end_time,
	event_date,
	number_of_guest,
	per_person_price,
	first_name,
	last_name,
	ticket_remaining
	   
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

	//	fmt.Print(et)
	return event, nil
}

const createEventQuery = `
	INSERT INTO events(
		event_name,
		number_of_guest,
		per_person_price,
		event_date,
		event_start_time,
		event_end_time,
		event_type_id,
		speakers_id,
		ticket_remaining
	)
	VALUES(
		:event_name,
		:number_of_guest,
		:per_person_price,
		:event_date,
		:event_start_time,
		:event_end_time,
		:event_type_id,
		:speakers_id,
		:ticket_remaining
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

const selectByIdQuery = `
SELECT 
	events.id, 
	event_name,
	event_type_id,
	event_type_name,
	event_start_time,
	event_end_time,
	event_date,
	number_of_guest,
	per_person_price,
	speakers_id,
	first_name,
	last_name,
	ticket_remaining,
	status
	
FROM events

JOIN speakers ON 
	events.speakers_id = speakers.id

	JOIN event_type ON 
	events.event_type_id = event_type.id

WHERE events.id = $1
`

func (s *Storage) GetDataById(id string) (storage.Events, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Unable to parse String to integer")
	}
	jason := storage.Events{}
	err = s.db.Get(&jason, selectByIdQuery, i)
	return jason, err
}

const countEvent = `SELECT COUNT(id) FROM events`

func (s *Storage) CountEvent() int32 {
	var count int32

	err := s.db.QueryRow(countEvent).Scan(&count)
	if err != nil {
		log.Println("Unable to ge data")
	}
	return count
}

const dec = `
UPDATE 
	events
SET 
	ticket_remaining = ticket_remaining - $1
WHERE 	
	ticket_remaining > 0 AND id = $2;
`

func (s *Storage) DecrementRemainingTicketById(id, number_of_ticket int32) (storage.Events, error) {
	jason := storage.Events{}
	err := s.db.Get(&jason, dec, number_of_ticket, id)
	return jason, err
}
