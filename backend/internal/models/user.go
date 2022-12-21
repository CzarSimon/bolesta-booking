package models

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/timeutil"
)

// User inidvidual that can book cabins
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Salt      string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u User) String() string {
	return fmt.Sprintf("User(id=%s, createdAt=%v, updatedAt=%v)", u.ID, u.CreatedAt, u.UpdatedAt)
}

// User inidvidual that can book cabins
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if r.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if len(r.Password) < 8 {
		return fmt.Errorf("password must be longer than 8 characters")
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

func (r CreateUserRequest) User(password, salt string) User {
	now := timeutil.Now()

	return User{
		ID:        id.New(),
		Name:      r.Name,
		Email:     r.Email,
		Password:  password,
		Salt:      salt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
