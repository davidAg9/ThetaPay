package services

import (
	"os"

	"github.com/keep94/toolbox/mailer"
)

func SendTransactionalEmail() (bool, error) {
	return true, nil
}

//send email to a sign user
func SendSignUpEmail(userMail string) {
	var email mailer.Email
	pass := os.Getenv("MAILPASS")
	myMailer := mailer.New("shane.qoubby@gmail.com", pass)
	email.To = []string{userMail}
	email.Subject = "Theta Email Verification Service"
	//TODO:CREATE EMAIL TEMPLATE
	emailTempLate := `
	<h1>Click the link below to verify your email</h1>
	<br>
	<p>http://localhost:8080/customer/email/[token]</p>
	`
	email.Body = emailTempLate

	myMailer.Send(email)

}
