package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	Manager  UserRole = "manager"
	Attendee UserRole = "attendee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"` // The `json:"-"` tag is used to exclude the password field from the JSON response.
	Role      UserRole  `json:"role" gorm:"text;default:attendee"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) AfterCreate(db *gorm.DB) (err error) {
	if u.ID == 1 {
		db.Model(u).Update("role", Manager)
	}
	return
}
