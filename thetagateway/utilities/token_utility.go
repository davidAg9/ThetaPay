package utilities

import (
	"crypto/sha256"
	"fmt"
	"log"
	mRand "math/rand"
	"os"
	"strings"
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

func GenerateApiKey(word string) (hash *string, key *string, err error) {
	hasher := sha256.New()
	_, err = hasher.Write([]byte(word))
	if err != nil {
		return nil, nil, err
	}
	data := hasher.Sum(nil)
	hashWord := fmt.Sprintf("%x", data)
	return &hashWord, &word, nil
}

func GenerateRandomString() (newWord *string, err error) {
	mRand.Seed(time.Now().Unix())

	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 32
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := mRand.Intn(len(charSet))
		randomChar := charSet[random]
		_, err := output.WriteString(string(randomChar))
		if err != nil {
			return nil, err
		}

	}
	word := output.String()
	return &word, nil

}

func ValidateAPIToken(key *string, secret *string) (valid bool, err error) {
	hasher := sha256.New()
	_, err = hasher.Write([]byte(*secret))
	if err != nil {
		return false, err
	}
	sum := hasher.Sum(nil)

	hash := fmt.Sprintf("%x", sum)

	if *key == hash {
		return true, nil
	}

	return false, nil

}
