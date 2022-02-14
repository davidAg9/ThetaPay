package utilities

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ThetaCredentials struct {
	Email    *string
	FullName *string
	Uid      *string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, fullName string, uid string) (signedToken string, err error) {
	claims := &ThetaCredentials{
		Email:    &email,
		FullName: &fullName,
		Uid:      &uid,
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
