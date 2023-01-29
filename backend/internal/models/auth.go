package models

import "fmt"

const UserRole = "USER"

// LoginRequest request to authenticate a user and recieve an auth token.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) String() string {
	return fmt.Sprintf("LoginRequest(email=%s)", l.Email)
}

// AuthenticatedResponse request to authenticate a user and recieve an auth token.
type AuthenticatedResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (r AuthenticatedResponse) String() string {
	return fmt.Sprintf("AuthenticatedResponse(user=%s, token=%s)", r.User, r.Token)
}
