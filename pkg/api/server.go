package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "readon/cmd/api/docs"
	handler "readon/pkg/api/handler"

	middleware "readon/pkg/api/middleware"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, productHandler *handler.ProductHandler, adminHandler *handler.AdminHandler) *ServerHTTP { //
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	// Swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//SIGN UP
	engine.GET("/signup", userHandler.GetSignup)
	engine.POST("/signup", userHandler.SaveUser)

	//general login
	engine.GET("/login", middleware.AuthorizationMiddleware, userHandler.UserHome)
	engine.POST("/login", userHandler.UserLogin)

	//  admin login
	engine.GET("/adminlogin", middleware.AuthorizationMiddleware)
	engine.POST("/adminlogin", adminHandler.Login)

	// Auth middleware
	users := engine.Group("/user", middleware.AuthorizationMiddleware)
	admin := engine.Group("/admin", middleware.AuthorizationMiddleware)

	//users.GET("users", userHandler.FindAll)
	users.GET("/search", productHandler.ListProducts)

	admin.DELETE("user/:id", userHandler.Delete)
	admin.GET("user/:id", userHandler.FindByID)
	admin.GET("/users", adminHandler.ListUser)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
