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
