package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	User struct {
		ID        int32     `db:"id"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Password  string    `db:"password"`
		IsActive  bool      `db:"is_active"`
		IsAdmin   bool      `db:"is_admin"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	EventType struct {
		ID            int32     `db:"id"`
		EventTypeName string    `db:"event_type_name"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}

	Speakers struct {
		ID        int32     `db:"id"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Phone     string    `db:"phone"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Address   string    `db:"address"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	Events struct {
		ID               int32     `db:"id"`
		EventTypeId      int32     `db:"event_type_id"`
		EventTypeName    string    `db:"event_type_name"`
		EventName        string    `db:"event_name"`
		NumberOfGuest    int32     `db:"number_of_guest"`
		EventStartTime   string    `db:"event_start_time"`
		EventEndTime     string    `db:"event_end_time"`
		EventDate        string    `db:"event_date"`
		PerPersonPrice   int32     `db:"per_person_price"`
		SpeakerId        int32     `db:"speakers_id"`
		SpeakerFirstName string    `db:"first_name"`
		SpeakerLastName  string    `db:"last_name"`
		Status           bool      `db:"status"`
		CreatedAt        time.Time `db:"created_at"`
		UpdatedAt        time.Time `db:"updated_at"`
	}

	Feedback struct {
		ID        int32     `db:"id"`
		UserId    int32     `db:"user_id"`
		Email     string    `db:"email"`
		Username  string    `db:"username"`
		Message   string    `db:"message"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

const nameLength = "Name should be 4 to 30 Characters"
const firstNameRequired = "First Name is Required"
const lastNameRequired = "Last Name is Required"
const emailIsRequired = "Email is Required"
const passwordIsRequired = "Email is Required"
const usernameIsRequired = "User name is Required"
const eventTypeNameIsRequired = "Event Type name is Required"

// User validation
func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName,
			validation.Required.Error(firstNameRequired),
			validation.Length(5, 30).Error(nameLength),
		),

		validation.Field(&a.LastName,
			validation.Required.Error(lastNameRequired),
			validation.Length(5, 30).Error(nameLength),
		),

		validation.Field(&a.Email,
			validation.Required.Error(emailIsRequired),
			validation.Length(5, 30).Error(nameLength),
		),

		validation.Field(&a.Username,
			validation.Required.Error(usernameIsRequired),
			validation.Length(5, 30).Error(nameLength),
		),
		validation.Field(&a.Password,
			validation.Required.Error(passwordIsRequired),
			validation.Length(5, 30).Error(nameLength),
		),
	)
}

// Event Type Validation
func (a EventType) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.EventTypeName,
			validation.Required.Error(eventTypeNameIsRequired),
			validation.Length(5, 30).Error(nameLength),
		),
	)
}
