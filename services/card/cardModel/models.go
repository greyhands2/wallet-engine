package cardModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Card_id            string             `json:"card_id"`
	Authorization_code string             `json:"authorization_code" validate:"required"`
	Card_type          string             `json:"card_type" validate:"required"`
	Last4              string             `json:"last4" validate:"required"`
	Exp_month          string             `json:"exp_month" validate:"required"`
	Exp_year           string             `json:"exp_year" validate:"required"`
	Card               string             `json:"card" validate:"required"`
	Channel            string             `json:"channel" validate:"required"`
	Signature          string             `json:"signature" validate:"required"`
	Reusable           bool               `json:"reusable" validate:"required"`
	Country_code       string             `json:"country_code" validate:"required"`
	User_id            string             `json:"user_id" validate:"required"`
	Created_at         time.Time          `json:"created_at"`
	Updated_at         time.Time          `json:"updated_at"`
}
