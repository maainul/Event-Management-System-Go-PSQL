package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
	"log"
	"strconv"
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

const getSpeakerById = `
	Select * 
	FROM speakers
	WHERE id = $1
`

func (s *Storage) GetSpeakerById(id string) (storage.Speakers, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Unable to parse String to integer")
	}
	jason := storage.Speakers{}
	err = s.db.Get(&jason, getSpeakerById, i)
	return jason, err
}

const update = `
	UPDATE 
		speakers 
	SET 

		first_name=:first_name,
		last_name =:last_name,
		phone =:phone,
		username=: username,
		email=: email,
		address=:address

	WHERE 
		id =: id`

func (s *Storage) UpdateSpeaker(speaker storage.Speakers) {
	/* stmt, err := s.db.PrepareNamed(update)
	if err != nil {
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, speaker); err != nil {
		return 0, err
	}
	return id, nil */

	res, err := s.db.NamedExec(update, speaker)
	if err != nil {
		log.Println("Unable to update data")
	}
	fmt.Println("------>", res)

}
