package helpers

import (
	"errors"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
	uuid "github.com/satori/go.uuid"
)

var emailjetApiKey, emailjetSecretKey string

func SetEmailConfig(apiKey, secretKey string) error {
	if apiKey == "" {
		return errors.New("email servise api key is empty")
	}
	if secretKey == "" {
		return errors.New("email servise secret key is empty")
	}
	fmt.Println("sennf grid api key ", apiKey)
	emailjetApiKey = apiKey
	emailjetSecretKey = secretKey
	return nil
}

func GenerateAndSendOpt(email string) (string, error) {
	otpUUID := uuid.NewV4()
	otp := otpUUID.String()[:6]

	// from := mail.NewEmail("READON", "ahmdkabeerm@gmail.com")
	// subject := otp
	// to := mail.NewEmail(email, email)
	// plainTextContent := fmt.Sprintf("Your OTP is: %s", otp)
	// htmlContent := fmt.Sprintf("<p>Your OTP is: <strong>%s</strong></p>", otp)

	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	mailjetClient := mailjet.NewMailjetClient(emailjetApiKey, emailjetSecretKey)
	fmt.Println("email :", email)
	recipientName := "add RECIEPINT NAME"
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "ahmdkabeerm@gmail.com",
				Name:  "ReadOn Support",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  recipientName,
				},
			},
			Subject:  "Your OTP for ReadOn Account Verification",
			TextPart: fmt.Sprintf("Hello %s,\n\nYour OTP for verifying your ReadOn account is: %s.\nPlease enter this code on the verification page to complete the process.\n\nIf you did not request this, please ignore this email.\n\nThank you for using ReadOn!\n\nBest Regards,\nThe ReadOn Team", recipientName, otp),
			HTMLPart: fmt.Sprintf(`
                <h2>Hello %s,</h2>
                <p>Your OTP for verifying your ReadOn account is: <strong>%s</strong>.</p>
                <p>Please enter this code on the verification page to complete the process.</p>
                <p>If you did not request this, please ignore this email.</p>
                <br>
                <p>Thank you for using <strong>ReadOn</strong>!</p>
                <p>Best Regards,</p>
                <p>The ReadOn Team</p>`, recipientName, otp),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Data: %+v\n", res)
	}
	return otp, nil

}
