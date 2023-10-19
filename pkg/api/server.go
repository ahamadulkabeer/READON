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

func NewServerHTTP(userHandler *handler.UserHandler, productHandler *handler.ProductHandler, adminHandler *handler.AdminHandler, categoryHandler *handler.CategoryHandler) *ServerHTTP {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	// Swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//engine.LoadHTMLGlob("pkg/templates/*.html")

	// Auth middleware
	users := engine.Group("/user", middleware.UserAuthorizationMiddleware)
	admin := engine.Group("/admin", middleware.AdminAuthorizationMiddleware)
	category := admin.Group("/category")

	//user handlers
	engine.GET("/signup", userHandler.GetSignup)
	engine.POST("/signup", userHandler.SaveUser)
	engine.POST("/login", userHandler.UserLogin)
	engine.GET("/login", userHandler.GetLogin)
	engine.GET("/otplogin", userHandler.GetOtpLogin)
	engine.POST("/otplogin", userHandler.VerifyAndSendOtp)
	engine.POST("/verifyotp", userHandler.VerifyOtp)

	users.GET("/home", userHandler.UserHome, productHandler.ListProducts)
	users.GET("/books", productHandler.ListProducts)
	users.GET("/book/:id", productHandler.GetProduct)

	//admin handlers
	engine.GET("/adminlogin", adminHandler.GetLogin)
	engine.POST("/adminlogin", adminHandler.Login)

	admin.PUT("/user/:id", adminHandler.BlockOrUnBlock)
	admin.DELETE("/user/:id", adminHandler.Delete)
	admin.GET("/user/:id", adminHandler.FindByID)
	admin.GET("/users", adminHandler.ListUser)
	admin.POST("/addproduct", productHandler.Addproduct)

	//categories
	category.GET("/categorylist", categoryHandler.ListCategories)
	category.POST("/addcategory", categoryHandler.AddCategory)
	category.PUT("/updatecategory/:id", categoryHandler.UpdateCategory)
	category.DELETE("/deletecategory/:id", categoryHandler.DeleteCategory)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
