package types

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type HttpSuccessResponse struct {
	Message string `json:"message"`
}

type SingleFileUploadedSuccessResponse struct {
	FileName string `json:"fileName"`
	Message  string `json:"message"`
}

type MultipleFileUploadedSuccessResponse struct {
	FileNames []string `json:"fileNames"`
	Message   string   `json:"message"`
}

type DeleteSuccessResponse struct {
	FileName string `json:"fileName"`
	Message  string `json:"message"`
}

type User struct {
	ID        string `json:"id"`
	UserName  string `json:"username" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Verified  bool   `json:"verified"`
}

type User_Mongo struct {
	ID                string    `bson:"id"`
	UserName          string    `bson:"username, unique"`
	FirstName         string    `bson:"firstName"`
	LastName          string    `bson:"lastName"`
	Email             string    `bson:"email, unique"`
	Password          string    `bson:"password"`
	ProfilePicture    string    `bson:"profilePicture"`
	Verified          bool      `bson:"verified"`
	CreatedAt         time.Time `bson:"createdAt"`
	UpdatedAt         time.Time `bson:"updatedAt"`
	VerificationToken int       `bson:"verificationToken"`
}

type User_Mongo_Oauth struct {
	ID             string    `bson:"id"`
	UserName       string    `bson:"username, unique"`
	FirstName      string    `bson:"firstName"`
	LastName       string    `bson:"lastName"`
	Email          string    `bson:"email, unique"`
	ProfilePicture string    `bson:"profilePicture"`
	Verified       bool      `bson:"verified"`
	CreatedAt      time.Time `bson:"createdAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
}

type LogInDetails struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
