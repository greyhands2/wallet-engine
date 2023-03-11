package debitControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/config"
	aDebit "github.com/greyhands2/wallet-engine/services/debit/debitModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var debitCollection *mongo.Collection = config.OpenCollection(config.Client, "debit")

func GetDebits(reqRes *fiber.Ctx) error {
	//get user id
	user_id := reqRes.Locals("user_id").(string)

	cursor, err := debitCollection.Find(reqRes.Context(), bson.M{"user_id": user_id})

	if err != nil {
		return reqRes.Status(500).SendString(err.Error())
	}

	var debits []aDebit.Debit = make([]aDebit.Debit, 0)

	if err = cursor.All(reqRes.Context(), &debits); err != nil {

		return reqRes.Status(500).SendString(err.Error())
	}

	defer cursor.Close(reqRes.Context())
	return reqRes.Status(200).JSON(debits)

}

func GetDebit(reqRes *fiber.Ctx) error {
	//initiate credit struct
	var debit *aDebit.Debit
	//get user id
	user_id := reqRes.Locals("user_id").(string)

	//get creditId from params
	debit_id := reqRes.Params("debitId")

	//not make the query
	query := bson.M{
		"$and": []bson.M{
			bson.M{"debit_id": debit_id},
			bson.M{"user_id": user_id},
		},
	}

	err := debitCollection.FindOne(reqRes.Context(), query).Decode(&debit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return reqRes.Status(404).SendString("credit trail not found")
		}
		return reqRes.Status(500).SendString("Oopss!! something went wrong, please try again later")
	}

	return reqRes.Status(200).JSON(debit)
}
