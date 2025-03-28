package models

import "net/http"

type UserLoginRequest struct {
	UserEmail    string `json:"user_email" validate:"required,email,max=50"`
	UserPassword string `json:"user_password" validate:"required,max=50"`
}

type User struct {
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name" validate:"required,max=50"`
	UserEmail    string `json:"user_email" validate:"required,email,max=50"`
	UserPassword string `json:"user_password" validate:"required,max=50"`
	UserPhone    string `json:"user_phone" validate:"required,max=50"`
}

type UserInterface interface {
	Login(*http.Request) (string, error)
	Register(*http.Request) (string, error)
	DeleteUser(*http.Request) error
}
