package utilities

import (
	"log"

	"github.com/davidAg9/thetagateway/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CheckRoleType(g *gin.Context, role models.Role) error {
	return nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "password is incorrect"
		check = false
	}
	return check, msg
}
