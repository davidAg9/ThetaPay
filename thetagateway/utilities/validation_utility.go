package utilities

import (
	"log"
	"net/mail"
)

func IsEmailValid(email string) bool {
	log.Print(email)
	_, err := mail.ParseAddress(email)
	return err == nil
}
