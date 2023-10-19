package interfaces

import (
	"readon/pkg/models"
)

type AdminRepository interface {
	Login(models.Userlogindata) (int, bool)
	ListAdmins() ([]models.Admin, error)
}
