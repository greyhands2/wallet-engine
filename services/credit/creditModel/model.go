package creditModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credit struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Credit_id string             `json:"credit_id"`
	User_id   string             `json:"user_id"`
	Amount    int32              `json:"amount" validate:"required"`
	Status    string             `json:"status"`

	Paid bool `json:"paid"`

	Disabled bool   `json:"disabled"`
	Txn_ref  string `json:"txn_ref"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
