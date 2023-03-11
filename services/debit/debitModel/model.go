package debitModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Debit struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Debit_id  string             `json:"debit_id"`
	User_id   string             `json:"user_id"`
	Wallet_id string             `json:"wallet_id"`
	Amount    int32              `json:"amount" validate:"required"`
	Purpose   string             `json:"purpose" validate:"required"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
