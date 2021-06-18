package postgres

import "Event-Management-System-Go-PSQL/storage"

const getUser = `
	SELECT id, first_name, last_name, username, email from users
	WHERE id = $1
`

func (s *Storage) GetUser(id int32) (*storage.User, error) {
	user := storage.User{}
	if err := s.db.Get(&user, getUser, id); err != nil {
		return nil, err
	}
	return &user, nil
}
