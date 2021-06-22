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
