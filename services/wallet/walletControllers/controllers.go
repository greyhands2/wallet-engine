package walletControllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aCard "github.com/greyhands2/wallet-engine/services/card/cardModel"
	aCredit "github.com/greyhands2/wallet-engine/services/credit/creditModel"
	aDebit "github.com/greyhands2/wallet-engine/services/debit/debitModel"
	aWallet "github.com/greyhands2/wallet-engine/services/wallet/walletModel"
	"github.com/greyhands2/wallet-engine/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var walletCollection *mongo.Collection = config.OpenCollection(config.Client, "wallet")

var cardCollection *mongo.Collection = config.OpenCollection(config.Client, "card")
var creditCollection *mongo.Collection = config.OpenCollection(config.Client, "credit")
var debitCollection *mongo.Collection = config.OpenCollection(config.Client, "debit")
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
	var user_id string = reqRes.Locals("user_id").(string)

	var wallet *aWallet.Wallet

	err := walletCollection.FindOne(reqRes.Context(), bson.M{"user_id": user_id, "activated": true}).Decode(&wallet)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("wallet not found")
		}
		return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
	}

	return reqRes.Status(200).JSON(wallet)

}

func ChangeWalletStatus(reqRes *fiber.Ctx) error {
	var user_id string = reqRes.Locals("user_id").(string)
	var status_type string = reqRes.Locals("status_type").(string)

	query := bson.M{"user_id": user_id}
	var updateValue bool
	var resultMessage string
	if status_type == "deactivate" {
		updateValue = false
		resultMessage = "Successfully deactivated wallet"
	} else {
		updateValue = true
		resultMessage = "Successfully activated wallet"
	}
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{"$set": bson.M{"activated": updateValue, "updated_at": updated_at}}

	err := walletCollection.FindOneAndUpdate(reqRes.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("wallet not found")
		}
		return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
	}

	return reqRes.Status(200).SendString(resultMessage)
}

func CreditWallet(reqRes *fiber.Ctx) error {
	var card *aCard.Card
	var credit *aCredit.Credit
	//retrieve cardId
	card_id := reqRes.Params("cardId")

	//retrieve user_id
	user_id := reqRes.Locals("user_id").(string)

	query := bson.M{
		"$and": []bson.M{
			bson.M{"card_id": card_id},
			bson.M{"user_id": user_id},
		},
	}

	// now let's fetch card details
	err := cardCollection.FindOne(reqRes.Context(), query).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("Card not found")
		}
		return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
	}

	//get amount from request body
	if err = reqRes.BodyParser(&credit); err != nil {
		return reqRes.Status(400).SendString(err.Error())
	}

	//validate the struct for just amount
	validationError := validate.Struct(credit)
	if validationError != nil {
		return reqRes.Status(400).SendString(validationError.Error())
	}

	//generate transaction reference
	txn_ref := utils.RandomString(15)

	credit.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	credit.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	credit.ID = primitive.NewObjectID()
	credit.Credit_id = credit.ID.Hex()
	credit.User_id = user_id
	credit.Status = "pending"

	credit.Paid = false

	credit.Disabled = false
	credit.Txn_ref = txn_ref

	//not lets create a credit trail

	_, insetErr := creditCollection.InsertOne(reqRes.Context(), credit)
	//check for error in insertion of card
	if insetErr != nil {
		return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
	}
	//now we have created a transaction trail
	//it is expected here that at this point that the credit  and card struct details respectively are to be sent to the payment gateway for example paystack which would initiate a transaction and when it is completed the details of the completion would hit a webhook url we provide them , that url would hit a function in our system in the credit controllers to complete the transaction and then reflect it into the wallet balance

	//but for the sake of testing here i would update the wallet and the credit trail as successful

	//initiate credit successful
	query = bson.M{
		"$and": []bson.M{
			bson.M{"credit_id": credit.Credit_id},
			bson.M{"user_id": credit.User_id},
		},
	}

	update := bson.M{
		"$set": bson.M{"status": "successful",
			"paid": true},
	}
	err = creditCollection.FindOneAndUpdate(reqRes.Context(), query, update).Err()

	if err != nil {
		return reqRes.Status(500).SendString("Something went wrong please contact customer care")
	}

	//now let's update the wallet
	query = bson.M{"user_id": user_id}

	update = bson.M{
		"$inc": bson.M{"balance": credit.Amount},
	}

	err = walletCollection.FindOneAndUpdate(reqRes.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("Wallet not found")
		}
		return reqRes.Status(500).SendString("Something went wrong please contact customer care")
	}

	return reqRes.Status(200).SendString("Credit Successful")

}

func DebitWallet(reqRes *fiber.Ctx) error {
	//get user id... the user_id
	user_id := reqRes.Locals("user_id").(string)
	//get wallet_id from params
	wallet_id := reqRes.Params("walletId")
	//set debit struct
	var debit *aDebit.Debit

	//get amount from request body
	if err := reqRes.BodyParser(&debit); err != nil {
		return reqRes.Status(400).SendString(err.Error())
	}

	//validate amount and purpose in debit struct
	validationError := validate.Struct(debit)
	if validationError != nil {
		return reqRes.Status(400).SendString(validationError.Error())
	}

	debit.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	debit.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	debit.ID = primitive.NewObjectID()
	debit.Debit_id = debit.ID.Hex()
	debit.User_id = user_id
	debit.Wallet_id = wallet_id

	//now let's substract the amount from the wallet first then on success we create a debit transaction

	query := bson.M{
		"$and": []bson.M{
			bson.M{"user_id": user_id},
			bson.M{"wallet_id": wallet_id},
		},
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": -debit.Amount,
		},
	}

	updateErr := walletCollection.FindOneAndUpdate(reqRes.Context(), query, update).Err()

	if updateErr != nil {
		if updateErr == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("Wallet not found")
		}
		return reqRes.Status(500).SendString("Oooopss!!! Something went wrong, please try again later")
	}

	//now let's create the debit transaciton trail, it would be ideal to put this in a worker queur or goroutine because if there's a server
	_, debitErr := debitCollection.InsertOne(reqRes.Context(), debit)

	if debitErr != nil {
		return reqRes.Status(500).SendString("Something went wrong please contact customer care")
	}

	return reqRes.Status(200).SendString("Debit transaction successful")

}
