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

	mailjetClient := mailjet.NewMailjetClient(emailjetApiKey, emailjetSecretKey)
	fmt.Println("email :", email)
	recipientName := "add User"
	htmlContent := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 20px auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px;">
			<h2 style="color: #333;">Welcome to ReadOn!</h2>
			<p>Hello %s,</p>
			<p>Thank you for signing up at <strong>ReadOn</strong>. Please use the following One-Time Password (OTP) to verify your email address:</p>
			<div style="text-align: center; margin: 20px 0;">
				<p style="font-size: 24px; font-weight: bold; color: #4CAF50;">%s</p>
			</div>
			<p>If you did not request this, please ignore this email.</p>
			<p>Thanks,<br/>The ReadOn Team</p>
			<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
			<p style="font-size: 12px; color: #888;">This email was sent by ReadOn. If you have any questions, please contact us at support@readon.com.</p>
		</div>
	</body>
	</html>
`, recipientName, otp)
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
			HTMLPart: htmlContent,
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
