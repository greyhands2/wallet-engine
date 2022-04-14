package cardControllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aCard "github.com/greyhands2/wallet-engine/services/card/cardModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cardCollection *mongo.Collection = config.OpenCollection(config.Client, "card")
var validate = validator.New()

func AddCard(reqRes *fiber.Ctx) error {

	var card *aCard.Card
	if err := reqRes.BodyParser(&card); err != nil {
		return reqRes.Status(400).SendString(err.Error())
	}
	//get user id
	var user_id string = reqRes.Locals("user_id").(string)
	//set user_id
	card.User_id = user_id
	//validate card struct
	validationError := validate.Struct(card)
	if validationError != nil {
		return reqRes.Status(400).SendString(validationError.Error())
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{"user_id": user_id},
			bson.M{"authorization_code": card.Authorization_code},
			bson.M{"signature": card.Signature},
		},
	}
	//check if exact card alread exists before inserting
	err := cardCollection.FindOne(reqRes.Context(), query).Err()

	if err != nil {
		//if there is no document it means we can add the card
		if err == mongo.ErrNoDocuments {
			card.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			card.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			card.ID = primitive.NewObjectID()
			card.Card_id = card.ID.Hex()

			//before inserting the card it is assumed that some sort of card verification is to be donw with an external payment gateway service e.g paystack before proceeding to save the card in our database

			// insert card
			cardRes, err := cardCollection.InsertOne(reqRes.Context(), card)
			//check for error in insertion of card
			if err != nil {
				return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
			}

			return reqRes.Status(200).JSON(cardRes)

		} else {
			// else then it is a server error
			return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
		}

	}

	return reqRes.Status(409).SendString("You already have this card added")
}

func GetCards(reqRes *fiber.Ctx) error {

	var user_id string = reqRes.Locals("user_id").(string)

	cursor, err := cardCollection.Find(reqRes.Context(), bson.M{"user_id": user_id})

	if err != nil {
		return reqRes.Status(500).SendString(err.Error())
	}

	var cards []aCard.Card = make([]aCard.Card, 0)

	if err = cursor.All(reqRes.Context(), &cards); err != nil {

		return reqRes.Status(500).SendString(err.Error())
	}

	defer cursor.Close(reqRes.Context())
	return reqRes.Status(200).JSON(cards)
}

func GetCard(reqRes *fiber.Ctx) error {

	var card *aCard.Card
	var user_id string = reqRes.Locals("user_id").(string)

	cardId, err := primitive.ObjectIDFromHex(reqRes.Params("cardId"))
	if err != nil {
		return reqRes.SendStatus(400)
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{"user_id": user_id},
			bson.M{"_id": cardId},
		},
	}
	err = cardCollection.FindOne(reqRes.Context(), query).Decode(&card)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("Card not found")
		}
		return reqRes.Status(500).SendString("Oopss!! Something went wrong, please try again later")
	}

	return reqRes.Status(200).JSON(card)
}

func RemoveCard(reqRes *fiber.Ctx) error {
	var user_id string = reqRes.Locals("user_id").(string)

	idParam := reqRes.Params("cardId")

	cardId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return reqRes.SendStatus(400)
	}

	query := bson.M{
		"$and": []bson.M{
			bson.M{"user_id": user_id},
			bson.M{"_id": cardId},
		},
	}

	res, err := cardCollection.DeleteOne(reqRes.Context(), query)

	if err != nil {
		return reqRes.SendStatus(500)
	}

	if res.DeletedCount < 1 {
		return reqRes.SendStatus(404)
	}

	return reqRes.Status(200).JSON("Card Sucessfully deleted")

}
