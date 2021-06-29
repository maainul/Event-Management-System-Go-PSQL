package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
)

const sp = `SELECT * from speakers`

func (s *Storage) GetSpeakers() ([]storage.Speakers, error) {
	speakers_list := make([]storage.Speakers, 0)
	if err := s.db.Select(&speakers_list, sp); err != nil {
		return nil, err
	}
	// fmt.Print(et)
	return speakers_list, nil
}

const createSpeakerQuery = `
	INSERT INTO speakers(
		first_name,
		last_name,
		phone,
		username,
		email,
		address
	)
	VALUES(
		:first_name,
		:last_name,
		:phone,
		:username,
		:email,
		:address
	)
	RETURNING id
	`

func (s *Storage) CreateSpeaker(usr storage.Speakers) (int32, error) {
	stmt, err := s.db.PrepareNamed(createSpeakerQuery)
	if err != nil {
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, usr); err != nil {
		return 0, err
	}
	return id, nil
}
