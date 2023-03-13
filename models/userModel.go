package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name"  validate:"required,min=2,max=100"`
	Email        string             `json:"email" bson:"email" validate:"email,required"`
	Password     string             `json:"password" validate:"required,min=6,max=20"`
	UserType     string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Token        string             `json:"token"`
	RefreshToken string             `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
