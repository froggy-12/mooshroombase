package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PongResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type User_Mongo struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
