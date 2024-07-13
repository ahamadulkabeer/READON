package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	ListUsers(models.Pagination) ([]domain.User, int, error)
	FindByEmail(email string) (domain.User, error)
	FindByID(id uint) (domain.User, error)
	DeleteUser(user domain.User) error
	BlockOrUnBlock(int) (bool, error)

	//Authorise(user models.Userlogindata) (int, bool, error)
	CheckForEmail(email string) (bool, error)
	SaveOtp(string, string) error
	VerifyOtp(string, string) error
}
