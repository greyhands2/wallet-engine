package cardRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/middleware"
	card "github.com/greyhands2/wallet-engine/services/card/cardControllers"
)

var HandleCardRoutes = func(router fiber.Router) {
	router.Post("/", middleware.Protect, card.AddCard)
	router.Get("/", middleware.Protect, card.GetCards)
	router.Get("/:cardId", middleware.Protect, card.GetCard)
	router.Delete("/:cardId", middleware.Protect, card.RemoveCard)

}
