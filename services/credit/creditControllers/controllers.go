package creditControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aCredit "github.com/greyhands2/wallet-engine/services/credit/creditModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var creditCollection *mongo.Collection = config.OpenCollection(config.Client, "credit")

func GetCredits(reqRes *fiber.Ctx) error {
	//get user id
	user_id := reqRes.Locals("user_id").(string)

	cursor, err := creditCollection.Find(reqRes.Context(), bson.M{"user_id": user_id})

	if err != nil {
		return reqRes.Status(500).SendString(err.Error())
	}

	var credits []aCredit.Credit = make([]aCredit.Credit, 0)

	if err = cursor.All(reqRes.Context(), &credits); err != nil {

		return reqRes.Status(500).SendString(err.Error())
	}

	defer cursor.Close(reqRes.Context())
	return reqRes.Status(200).JSON(credits)

}

func GetCredit(reqRes *fiber.Ctx) error {
	//initiate credit struct
	var credit *aCredit.Credit
	//get user id
	user_id := reqRes.Locals("user_id").(string)

	//get creditId from params
	credit_id := reqRes.Params("creditId")

	//not make the query
	query := bson.M{
		"$and": []bson.M{
			bson.M{"credit_id": credit_id},
			bson.M{"user_id": user_id},
		},
	}

	err := creditCollection.FindOne(reqRes.Context(), query).Decode(&credit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("credit trail not found")
		}
		return reqRes.Status(500).SendString("Oopss!! something went wrong, please try again later")
	}

	return reqRes.Status(200).JSON(credit)
}
