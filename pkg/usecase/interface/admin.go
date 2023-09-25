package interfaces

import (
	"context"
	"readon/pkg/models"
)

type AdminUseCase interface {
	Login(context.Context, models.Userlogindata) (int, bool)
	ListUser(context.Context) ([]models.ListOfUser, error)
}
