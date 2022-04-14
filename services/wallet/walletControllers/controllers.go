package walletControllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aWallet "github.com/greyhands2/wallet-engine/services/wallet/walletModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var walletCollection *mongo.Collection = config.OpenCollection(config.Client, "wallet")

var validate = validator.New()

//create wallet is only used if for some reason the user wallet was unable to be created during sign up
func CreateWallet(reqRes *fiber.Ctx) error {
	var wallet *aWallet.Wallet
	Created_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	id := primitive.NewObjectID()
	wallet_id := id.Hex()
	var balance float32 = 0.00
	// we have to use type assertion here because values stored in reqRes.Locals are of type interface so we use type assertion to retrieve the real type and then assign them
	var username string = reqRes.Locals("username").(string)
	var user_id string = reqRes.Locals("user_id").(string)

	wallet = &aWallet.Wallet{Username: username, User_id: user_id, Created_at: Created_at, Updated_at: Updated_at, ID: id, Wallet_id: wallet_id, Activated: true, Balance: balance}

	count, err := walletCollection.CountDocuments(reqRes.Context(), bson.M{"user_id": user_id})

	if err != nil {
		return reqRes.Status(500).SendString("Oopss!!! Somehting went wrong, please try again later")
	}

	if count > 0 {
		return reqRes.Status(409).SendString("You already have a Wallet")
	}

	walletRes, err := walletCollection.InsertOne(reqRes.Context(), wallet)
	if err != nil {
		return reqRes.Status(500).SendString("Oopss!!! Somehting went wrong, please try again later")
	}
	return reqRes.Status(200).JSON(walletRes)

}

func GetWalletBalance(reqRes *fiber.Ctx) error {
	return nil
}
