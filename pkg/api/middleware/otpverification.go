package middleware

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func GenerateAndSendOpt(email string) (string, error) {
	otpUUID := uuid.NewV4()
	otp := otpUUID.String()[:6]

	from := mail.NewEmail("READON", "ahmdkabeerm@gmail.com")
	subject := otp
	to := mail.NewEmail(email, email)
	plainTextContent := fmt.Sprintf("Your OTP is: %s", otp)
	htmlContent := fmt.Sprintf("<p>Your OTP is: <strong>%s</strong></p>", otp)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	//os.Getenv("SENDGRID_API_KEY")
	client := sendgrid.NewSendClient("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")

	_, err := client.Send(message)
	if err != nil {
		return " ", err
	}
	return otp, nil

}

func VerifyOtp(email, otp string) bool {
	return false
}
