package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *fiber.Ctx) error {

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return c.Status(401).SendString("unauthorized")
	}

	claims, err := validateToken(token[0])

	if err != nil {
		return c.Status(401).SendString("unauthorized")
	}

	fmt.Printf("claims ->  %+v", claims)

	return c.Next()
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		fmt.Println(err)
	}

	fmt.Println(claims["userID"], claims["email"], claims["expiredAt"])

	// expiredAtFloat, ok := claims["expiredAt"].(float64)
	expiredAtStr, ok := claims["expiredAt"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse token")
	}

	// if time.Now().Unix() > int64(expiredAtFloat) {
	// 	return nil, fmt.Errorf("token expired")
	// }

	expiredAt, err := time.Parse(time.RFC3339, expiredAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token")
	}

	now := time.Now()
	fmt.Println("now : ", now, "expired : ", expiredAt)

	if expiredAt.Before(now) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
