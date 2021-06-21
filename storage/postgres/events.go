package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"

)

const e = `
SELECT events.id, event_name,event_type_name,start_time,end_time,event_date,number_of_guest,per_person_price,first_name, last_name 

FROM events 

JOIN speakers ON events.speakers_id = speakers.id

JOIN event_type ON events.event_type_id = event_type.id;
`

func (s *Storage) GetEvent() ([]storage.Events, error) {
	event := make([]storage.Events, 0)
	if err := s.db.Select(&event, e); err != nil {
		return nil, err
	}
	
	
	fmt.Print(et)
	return event, nil
}
