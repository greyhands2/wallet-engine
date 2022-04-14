package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/greyhands2/wallet-engine/utils"
)

func Protect(reqRes *fiber.Ctx) error {
	token := reqRes.Get("Authorization")

	splitToken := strings.Split(token, "Bearer ")
	clientToken := splitToken[1]

	if clientToken == "" {
		return reqRes.Status(401).SendString("Thread carefully, you roam on unauthorized waters")
	}

	claims, err := utils.ValidateToken(clientToken)

	if err != "" {
		return reqRes.Status(401).SendString("Thread carefully you roam on unauthorized waters")
	}

	reqRes.Locals("email", claims.Email)
	reqRes.Locals("username", claims.Username)

	reqRes.Locals("user_id", claims.Uid)

	reqRes.Next()
	return nil
}
