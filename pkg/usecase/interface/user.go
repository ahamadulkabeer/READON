package interfaces

import (
	"readon/pkg/api/responses"
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type UserUseCase interface {
	Save(userInput models.SignupData) responses.Response

	UpdateUser(user models.UserUpdateData) (domain.User, error)

	UserLogin(userinput models.LoginData) responses.Response

	GetUserProfile(int) (domain.User, error)

	DeleteUserAccount(int) error

	VerifyAndSendOtp(string) error

	VerifyOtp(string, string) error
}
