package entity

import "github.com/google/uuid"

type User struct {
	Id        uuid.UUID `col:"id" id:"Id"`
	LastName  string    `col:"last_name"`
	FirstName string    `col:"first_name"`
	Username  string    `col:"username"`
}
