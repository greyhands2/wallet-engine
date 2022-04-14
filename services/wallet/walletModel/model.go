package WalletModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Username string `json:"username" validate:"required,min=2,max=100"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    string    `json:"user_id"`
	Wallet_id  string    `json:"wallet_id"`
	Activated  bool      `json:"activated"`
	Balance    float32   `json:balance`
}
