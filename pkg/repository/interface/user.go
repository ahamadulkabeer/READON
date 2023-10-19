package interfaces

import (
	"readon/pkg/domain"
	"readon/pkg/models"
)

type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	ListUsers(models.Pagination, int) ([]domain.User, int, error)
	FindByEmail(email string) (domain.User, error)
	FindByID(id uint) (domain.User, error)
	DeleteUser(user domain.User) error
	BlockOrUnBlock(int) bool
	Authorise(user models.Userlogindata) (int, bool, error)
	CheckForEmail(string) error
	SaveOtp(string, string) error
	VerifyOtp(string, string) error
}
