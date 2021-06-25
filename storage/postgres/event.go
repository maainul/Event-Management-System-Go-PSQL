package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
)

const et = `SELECT * from event_type`

func (s *Storage) GetEventType() ([]storage.EventType, error) {
	event := make([]storage.EventType, 0)
	if err := s.db.Select(&event, et); err != nil {
		return nil, err
	}
	// fmt.Print(et)
	return event, nil
}

const createEvent = `
		INSERT INTO event_type(event_type_name)
		VALUES (:event_type_name)
		RETURNING id
		`

func (s *Storage) CreateEventType(event_type storage.EventType) (int32, error) {
	stmt, err := s.db.PrepareNamed(createEvent)
	if err != nil {
		return 0, err
	}

	var id int32
	if err := stmt.Get(&id, event_type); err != nil {
		return 0, err
	}
	return id, nil
}
