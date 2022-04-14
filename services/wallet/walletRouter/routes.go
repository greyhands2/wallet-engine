package walletRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/middleware"
	wallet "github.com/greyhands2/wallet-engine/services/wallet/walletControllers"
)

var HandleWalletRoutes = func(router fiber.Router) {
	//create wallet is only used if for some reason the user wallet was unable to be created during sign up
	router.Post("/:userId", middleware.Protect, wallet.CreateWallet)
	router.Get("/:userId", middleware.Protect, wallet.GetWalletBalance)
}
