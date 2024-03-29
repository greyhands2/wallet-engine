package main

import (
	"log"
	
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	user "github.com/greyhands2/wallet-engine/services/user/userRouter"

	wallet "github.com/greyhands2/wallet-engine/services/wallet/walletRouter"

	card "github.com/greyhands2/wallet-engine/services/card/cardRouter"
	credit "github.com/greyhands2/wallet-engine/services/credit/creditRouter"
	debit "github.com/greyhands2/wallet-engine/services/debit/debitRouter"
)





func main() {

	app := fiber.New()
	//handle panics if any
	app.Use(recover.New())
	api := app.Group("/api")

	user.HandleUserRoutes(api.Group("/user"))

	wallet.HandleWalletRoutes(api.Group("/wallet"))

	card.HandleCardRoutes(api.Group("/card"))
	credit.HandleCreditRoutes(api.Group("/credit"))
	debit.HandleDebitRoutes(api.Group("/debit"))
	log.Fatal(app.Listen(":3400"))

}
