package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
)

const createBookingQuery = `
INSERT INTO booking(
	event_id,
	user_id,
	number_of_ticket,
	total_amount
	)
	VALUES(
	  :event_id,
	  :user_id,
	  :number_of_ticket,
	  :total_amount
	 )
	RETURNING id
`

func (s *Storage) CreateBooking(booking storage.Booking) (int32, error) {
	stmt, err := s.db.PrepareNamed(createBookingQuery)
	if err != nil {
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, booking); err != nil {
		return 0, err
	}
	return id, nil
}
