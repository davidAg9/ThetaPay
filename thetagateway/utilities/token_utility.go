package utilities

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidAg9/thetagateway/models"
	jwt "github.com/dgrijalva/jwt-go"
)

type ThetaUserCredentials struct {
	Role     *models.Role
	UserName *string
	Uid      *string
	jwt.StandardClaims
}

type ThetaCustomerCredentials struct {
	Uid      *string
	FullName *string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

///Generate customer tokens to allow customer to get transactions , view and deposit into his account
func GenerateCustomerTokens(credentials ThetaCustomerCredentials) (signedToken string, err error) {

	credentials.ExpiresAt = time.Now().Local().Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, credentials).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateCustomerToken(signedToken string) (claims *ThetaCustomerCredentials, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&ThetaCustomerCredentials{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = fmt.Sprintf("the token is invalid ,Reason: %s", err.Error())
		return
	}

	claims, ok := token.Claims.(*ThetaCustomerCredentials)
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

func GenerateUserTokens(credentials *ThetaUserCredentials) (signedToken string, err error) {

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, credentials).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateUserToken(signedToken string) (claims *ThetaUserCredentials, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&ThetaUserCredentials{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = fmt.Sprintf("the token is invalid ,Reason: %s", err.Error())
		return
	}

	claims, ok := token.Claims.(*ThetaUserCredentials)
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

func GenerateApiKey(email *string) (signedToken string, err error) {
	//TODO:PROPER API KEY GENERATION
	// var dateCreated = time.Now().Local().Unix()
	creds := jwt.StandardClaims{

		Issuer: *email,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, creds).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateAPIToken(key string) (claims *jwt.StandardClaims, msg string) {
	token, err := jwt.ParseWithClaims(
		key,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = fmt.Sprintf("the api key is invalid ,Reason: %s", err.Error())
		return
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		msg = fmt.Sprintf("the key is invalid, Reason %s", err.Error())
		return
	}
	//TODO:PROPER API KEY VALIDATION

	return claims, msg
}
