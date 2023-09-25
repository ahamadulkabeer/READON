package interfaces

import (
	"context"
	"readon/pkg/models"
)

type AdminRepository interface {
	Login(context.Context, models.Userlogindata) (int, bool)
	ListUser(context.Context) ([]models.ListOfUser, error)
}
