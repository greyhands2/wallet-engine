package userRouter

import (
	"github.com/gofiber/fiber/v2"
	user "github.com/greyhands2/wallet-engine/services/user/userControllers"
)

var HandleUserRoutes = func(router fiber.Router) {
	router.Post("/signup", user.Signup)
	router.Post("/login", user.Login)

}
