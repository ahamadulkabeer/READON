package interfaces

import (
	"readon/pkg/api/responses"
	"readon/pkg/domain"
	"readon/pkg/models"
)

type UserUseCase interface {
	Save(userInput models.SignupData) responses.Response

	UpdateUser(user models.UserUpdateData) responses.Response

	UserLogin(userinput models.LoginData) responses.Response

	GetUserProfile(userID int) responses.Response

	DeleteUserAccount(userID int) responses.Response

	VerifyAndSendOtp(email string) responses.Response

	VerifyOtp(otpData domain.Otp) responses.Response
}
