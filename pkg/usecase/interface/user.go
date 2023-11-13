package interfaces

import (
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type UserUseCase interface {
	Save(models.SignupData) (domain.User, error)

	UpdateUser(user models.UpdateData) (domain.User, error)

	UserLogin(userinput models.Userlogindata) (int, bool, bool, error)

	GetUserProfile(int) (domain.User, error)

	DeleteUserAccount(int) error

	VerifyAndSendOtp(string) error

	VerifyOtp(string, string) error
}
