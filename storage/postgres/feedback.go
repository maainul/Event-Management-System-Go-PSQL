package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
)

const show = `SELECT feedback.id,username,email,message 

FROM feedback 

JOIN users ON feedback.user_id = users.id;`

func (s *Storage) GetFeedback() ([]storage.Feedback, error) {
	feedback_list := make([]storage.Feedback, 0)
	if err := s.db.Select(&feedback_list, show); err != nil {
		return nil, err
	}
	// fmt.Print(et)
	return feedback_list, nil
}

const createFeedback = `
		INSERT INTO feedback(
			message,
			user_id
		)
		VALUES (
			:message,
			:user_id
		)
		RETURNING id
		`

func (s *Storage) CreateFeedback(feedback storage.Feedback) (int32, error) {
	stmt, err := s.db.PrepareNamed(createFeedback)
	if err != nil {
		return 0, err
	}

	var id int32
	if err := stmt.Get(&id, feedback); err != nil {
		return 0, err
	}
	return id, nil
}
