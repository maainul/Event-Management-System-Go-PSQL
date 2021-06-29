package postgres

import (
	"Event-Management-System-Go-PSQL/storage"
	"fmt"
)

const getuserQuery = `
	SELECT id, first_name, last_name, username, email,password from users
`

func (s *Storage) GetUser() ([]storage.User, error) {
	user_list := make([]storage.User, 0)
	if err := s.db.Select(&user_list, getuserQuery); err != nil {
		return nil, err
	}
	return user_list, nil
}

const createUserQuery = `
	INSERT INTO users(
		first_name,
		last_name,
		username,
		email,
		password
	)
	VALUES(
		:first_name,
		:last_name,
		:username,
		:email,
		:password
	)
	RETURNING id
	`

func (s *Storage) CreateUser(usr storage.User) (int32, error) {
	stmt, err := s.db.PrepareNamed(createUserQuery)
	if err != nil {
		return 0, err
	}
	var id int32
	if err := stmt.Get(&id, usr); err != nil {
		return 0, err
	}
	return id, nil
}

/* const getUserEmailAndPass = `
	SELECT * FROM users WHERE email= $1 AND password= $2
`

func (s *Storage) GetUserEmailAndPass(email string, password string) (*storage.User, error) {
	jason := storage.User{}
	err := s.db.Get(&jason, getUserEmailAndPass, email, password)
	fmt.Print("Get email and pass  = ", jason)
	return &jason, err
}
*/
const getUserEmailAndPass = `
	SELECT * from users
	WHERE email = $1 AND password = $2
`

func (s *Storage) GetUserEmailAndPass(email, password string) (*storage.User, error) {
	user := storage.User{}
	if err := s.db.Get(&user, getUserEmailAndPass, email, password); err != nil {
		return nil, err
	}
	fmt.Print("Get email and pass  = ", user)
	return &user, nil
}
