package interfaces

import (
	"readon/pkg/models"
)

type AdminRepository interface {
	Login(models.LoginData) (int, bool)
	ListAdmins() ([]models.Admin, error)
}
