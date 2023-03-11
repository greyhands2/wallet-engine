package userModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Username        string    `json:"username" validate:"required,min=2,max=100"`
	Password        string    `json:"password" validate:"required,min=6""`
	PasswordConfirm string    `json:"passwordconfirm" validate:"required,min=6""`
	Email           string    `json:"email"  validate:"email,required"`
	Token           *string   `json:"token"`
	Refresh_token   *string   `json:"refresh_token"`
	Created_at      time.Time `json:"created_at"`
	Updated_at      time.Time `json:"updated_at"`
	User_id         string    `json:"user_id"`
}
