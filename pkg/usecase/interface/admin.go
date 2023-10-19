package interfaces

import (
	domain "readon/pkg/domain"
	"readon/pkg/models"
)

type AdminUseCase interface {
	Login(models.Userlogindata) (int, bool)
	ListAdmins() ([]models.Admin, error)
	ListUsers(*models.Pagination) ([]domain.User, int, error)
	FindByID(id uint) (domain.User, error)
	Delete(user domain.User) error
	BlockOrUnBlock(int) bool
}
