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
