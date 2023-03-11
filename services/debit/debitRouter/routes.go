package debitRouter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/middleware"
	debit "github.com/greyhands2/wallet-engine/services/debit/debitControllers"
)

var HandleDebitRoutes = func(router fiber.Router) {
	router.Get("/", middleware.Protect, debit.GetDebits)
	router.Get("/:debitId", middleware.Protect, debit.GetDebit)

}
