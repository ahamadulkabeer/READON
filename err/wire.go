//go:build wireinject
//+build wireinject

package di

import (
	http "readon/pkg/api"
	handler "readon/pkg/api/handler"
	config "readon/pkg/config"
	db "readon/pkg/db"
	repository "readon/pkg/repository"
	usecase "readon/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
