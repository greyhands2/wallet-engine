package walletRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/middleware"
	wallet "github.com/greyhands2/wallet-engine/services/wallet/walletControllers"
)

var HandleWalletRoutes = func(router fiber.Router) {
	//create wallet is only used if for some reason the user wallet was unable to be created during sign up
	//create wallet
	router.Post("/:userId", middleware.Protect, wallet.CreateWallet)
	//Get wallet balance
	router.Get("/:userId", middleware.Protect, wallet.GetWalletBalance)
	//Deactivate wallet
	router.Patch("/deactivate/:userId", middleware.Protect, func(reqRes *fiber.Ctx) error {
		reqRes.Locals("status_type", "deactivate")
		reqRes.Next()
		return nil
	}, wallet.ChangeWalletStatus)
	//Activate wallet
	router.Patch("/activate/:userId", middleware.Protect, func(reqRes *fiber.Ctx) error {
		reqRes.Locals("status_type", "activate")
		reqRes.Next()
		return nil
	}, wallet.ChangeWalletStatus)
}
