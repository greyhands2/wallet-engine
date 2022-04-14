package creditRouter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/middleware"
	credit "github.com/greyhands2/wallet-engine/services/credit/creditControllers"
)

var HandleCreditRoutes = func(router fiber.Router) {
	router.Get("/", middleware.Protect, credit.GetCredits)
	router.Get("/:creditId", middleware.Protect, credit.GetCredit)

}
