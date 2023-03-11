package userControllers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aUser "github.com/greyhands2/wallet-engine/services/user/userModel"
	aWallet "github.com/greyhands2/wallet-engine/services/wallet/walletModel"
	"github.com/greyhands2/wallet-engine/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "user")

var walletCollection *mongo.Collection = config.OpenCollection(config.Client, "wallet")
var validate = validator.New()

func Signup(reqRes *fiber.Ctx) error {
	//expecting user name password, password confirm and email

	var user aUser.User
	var err error

	// collect request body into user struct and handle error
	if err := reqRes.BodyParser(&user); err != nil {
		return reqRes.Status(400).SendString(err.Error())
	}

	//return error if both passwords are not the same
	if user.Password != user.PasswordConfirm {

		return reqRes.Status(409).SendString("Password do not match")
	}

	// lets validate the user struct
	structValidationError := validate.Struct(user)
	if structValidationError != nil {

		return reqRes.Status(400).SendString("Data Validation Error")
	}

	//ensure email is unique
	count, err := userCollection.CountDocuments(reqRes.Context(), bson.M{"email": user.Email})

	if count > 0 {

		return reqRes.Status(400).SendString("Username or Email already exist")
	}

	if err != nil {

		return reqRes.Status(400).SendString("Error occurred while validating email")

	}

	//ensure username is unique
	count, err = userCollection.CountDocuments(reqRes.Context(), bson.M{"username": user.Username})

	if err != nil {

		return reqRes.Status(400).SendString("Error occurred while validating username")

	}
	//send error if duplicate username or password exists
	if count > 0 {

		return reqRes.Status(400).SendString("Username or Email already exist")
	}

	//hash password
	password, hashErr := hashPassword(user.Password)

	if hashErr != nil {

		return reqRes.Status(500).SendString("Error processing data")
	}

	user.Password = password

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	//creating jwt tokens
	token, refreshToken, _ := utils.GenerateTokens(user.Email, user.Username, user.User_id)

	user.Token = &token
	user.Refresh_token = &refreshToken

	insertionNumber, insertionErr := userCollection.InsertOne(reqRes.Context(), user)

	if insertionErr != nil {
		return reqRes.Status(500).SendString(insertionErr.Error())
	}

	var wallet *aWallet.Wallet
	Created_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	id := primitive.NewObjectID()
	wallet_id := id.Hex()
	var balance float32 = 0.00
	wallet = &aWallet.Wallet{Username: user.Username, User_id: user.User_id, Created_at: Created_at, Updated_at: Updated_at, ID: id, Wallet_id: wallet_id, Activated: true, Balance: balance}
	//for microservice architecture replace this next step with send an event to event bus to get to the wallet service for initiation of wallet

	count, err = walletCollection.CountDocuments(reqRes.Context(), bson.M{"user_id": user.User_id})

	if err != nil {
		return reqRes.Status(500).SendString("Oopss!!! Somehting went wrong, please try again later")
	}

	if count > 0 {
		return reqRes.Status(409).SendString("You already have a Wallet")
	}
	_, insertionErr = walletCollection.InsertOne(reqRes.Context(), wallet)

	message := fmt.Sprintf("user %s created successfully", insertionNumber)

	if insertionErr != nil {
		walletErr := ", However your wallet creation" + "failed so you need to use the " +
			"Create Wallet button before you can use any of " + "our services"
		return reqRes.Status(500).SendString(message + walletErr)
	}

	return reqRes.Status(200).SendString(message)

}

func Login(reqRes *fiber.Ctx) error {

	var user aUser.User
	var foundUser aUser.User

	// collect request body into user struct and handle error
	if err := reqRes.BodyParser(&user); err != nil {
		return reqRes.Status(400).SendString(err.Error())
	}

	//find user and decode the bson result to the foundUser variable reference
	err := userCollection.FindOne(reqRes.Context(), bson.M{"email": user.Email}).Decode(&foundUser)

	if err != nil {
		return reqRes.Status(401).SendString("Incorrect email or password")
	}

	passwordIsValid, _ := validatePassword(foundUser.Password, user.Password)

	if passwordIsValid {
		return reqRes.Status(401).SendString("Incorrect email or password")
	}

	token, refreshToken, _ := utils.GenerateTokens(foundUser.Email, foundUser.Username, foundUser.User_id)

	utils.UpdateAllTokens(token, refreshToken, foundUser.User_id)

	foundUser.Password = ""
	foundUser.PasswordConfirm = ""

	return reqRes.Status(200).JSON(foundUser)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	var hashErr error
	if err != nil {
		hashErr = err
	}

	return string(hash), hashErr
}

func validatePassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""
	if err != nil {
		msg = "login or passowrd is incorrect"
		check = false
	}

	return check, msg

}
