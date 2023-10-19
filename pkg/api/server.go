package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "readon/cmd/api/docs"
	handler "readon/pkg/api/handler"
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
	users := engine.Group("/user")
	admin := engine.Group("/admin")
	//category := admin.Group("/category")

	//user handlers
	engine.GET("/signup", userHandler.GetSignup)
	engine.POST("/signup", userHandler.SaveUser)
	engine.POST("/login", userHandler.UserLogin)
	engine.GET("/login", userHandler.GetLogin)
	engine.GET("/otplogin", userHandler.GetOtpLogin)
	engine.POST("/otplogin", userHandler.VerifyAndSendOtp)
	engine.POST("/verifyotp", userHandler.VerifyOtp)

	users.DELETE("/account/:id", userHandler.DeleteUserAccount)
	users.GET("/profile/:id", userHandler.GetUserProfile)
	/*users.PUT("/editprofile",userHandler.EditProfile)
	users.GET("/readbook", userHandler.ReadBook)
	users.GET("/premium",userHandler.GetPremium)
	users.POST("/premium",userHandler.MakePremium)*/
	users.GET("/home", userHandler.UserHome, productHandler.ListProducts)
	users.GET("/books", productHandler.ListProducts)
	users.GET("/book/:id", productHandler.GetProduct)
	users.POST("/listbooks", productHandler.ListProductsForUSer)

	//admin handlers
	engine.GET("/adminlogin", adminHandler.GetLogin)
	engine.POST("/adminlogin", adminHandler.Login)

	admin.PUT("/blockuser/:id", adminHandler.BlockOrUnBlock)
	admin.DELETE("/user/:id", adminHandler.Delete)
	admin.GET("/user/:id", adminHandler.FindByID)
	admin.GET("/users", adminHandler.ListUsers)
	admin.GET("/admins", adminHandler.ListAdmins)
	admin.POST("/addproduct", productHandler.Addproduct)
	admin.DELETE("/deletebook/:id", productHandler.DeleteProduct)

	//categories
	admin.GET("/categorylist", categoryHandler.ListCategories)
	admin.POST("/addcategory", categoryHandler.AddCategory)
	admin.PUT("/updatecategory/:id", categoryHandler.UpdateCategory)
	admin.DELETE("/deletecategory/:id", categoryHandler.DeleteCategory)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
