package utilities

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidAg9/thetagateway/models"
	jwt "github.com/dgrijalva/jwt-go"
)

type ThetaCredentials struct {
	PhoneNumber *string
	FullName    *string
	Uid         *string
	Role        *models.Role
	jwt.StandardClaims
}

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
//TODO:SET SECRETE KEY
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(phoneNumber string, fullName string, role models.Role, uid string) (signedToken string, err error) {
	claims := &ThetaCredentials{
		PhoneNumber: &phoneNumber,
		FullName:    &fullName,
		Uid:         &uid,
		Role:        &role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateToken(signedToken string) (claims *ThetaCredentials, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&ThetaCredentials{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = fmt.Sprintf("the token is invalid ,Reason: %s", err.Error())
		return
	}

	claims, ok := token.Claims.(*ThetaCredentials)
	if !ok {
		msg = fmt.Sprintf("the token is invalid, Reason %s", err.Error())
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired %s", err.Error())
		return
	}
	return claims, msg
}
