package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "readon/cmd/api/docs"
	handler "readon/pkg/api/handler"
	"readon/pkg/api/middleware"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	adminHandler *handler.AdminHandler,
	categoryHandler *handler.CategoryHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHAndler,
	addressHandler *handler.AddressHandler,
	couponHandler *handler.CouponHandler) *ServerHTTP {

	engine := gin.New()
	engine.Use(cors.Default())

	//engine.LoadHTMLGlob("../templates/*")

	// Use logger from Gin
	engine.Use(gin.Logger())
	//to dowload invoice form browser without cookie authorisation :)
	engine.GET("/invoice/:orderId", orderHandler.DownloadInvoice)
	// Swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	engine.LoadHTMLGlob("pkg/templates/*.html")
	//user login / sign up
	engine.GET("/signup", userHandler.GetSignup)
	engine.POST("/signup", userHandler.SaveUser)
	engine.POST("/login", userHandler.UserLogin)
	engine.GET("/login", userHandler.GetLogin)
	engine.GET("/otplogin", userHandler.GetOtpLogin)
	engine.POST("/otplogin", userHandler.VerifyAndSendOtp)
	engine.POST("/verifyotp", userHandler.VerifyOtp)

	// admin login
	engine.GET("/adminlogin", adminHandler.GetLogin)
	engine.POST("/adminlogin", adminHandler.Login)

	//categories
	engine.GET("/categories", categoryHandler.ListCategories) //
	//book
	engine.GET("/books/:bookId", productHandler.GetProduct)  // ? not all book is getting ??
	engine.GET("/books", productHandler.ListProductsForUSer) //
	//home
	engine.GET("/home", userHandler.UserHome, productHandler.ListProducts) ///
	//web hook reciever (razor pay)
	engine.POST("/payment/verify", orderHandler.VerifyPayment) ///

	users := engine.Group("/user", middleware.UserAuthorizationMiddleware)
	{
		users.DELETE("/account", userHandler.DeleteUserAccount) //
		users.GET("/profile", userHandler.GetUserProfile)       //
		users.PUT("/profile", userHandler.UpdateUser)           ///

		// cart
		users.POST("/cart", cartHandler.AddToCart)        //
		users.GET("/cart", cartHandler.GetCart)           //
		users.PUT("/cart", cartHandler.UpdateQuantity)    //
		users.DELETE("/cart", cartHandler.DeleteFromCart) //
		// order
		users.POST("/orders", orderHandler.AddOrder)                  //
		users.DELETE("/orders/:orderId", orderHandler.CancelOrder)    //
		users.GET("/orders/:orderId", orderHandler.GetOrder)          //
		users.GET("/orders", orderHandler.ListOrders)                 //
		users.POST("/orders/:orderId/retry", orderHandler.RetryOrder) //
		// address
		users.POST("/addresses", addressHandler.AddAddress)                 //
		users.PUT("/addresses/:addressId", addressHandler.UpdateAddress)    //
		users.DELETE("/addresses/:addressId", addressHandler.DeleteAddress) //
		users.GET("/addresses/:addressId", addressHandler.GetAddress)       //
		users.GET("addresses", addressHandler.ListAddress)                  //
		//invoice
		users.GET("/invoice/:orderId", orderHandler.DownloadInvoice)
	}
	admin := engine.Group("/admin", middleware.AdminAuthorizationMiddleware)
	{
		//users
		admin.PUT("/users/:userId/block", adminHandler.BlockOrUnBlock) //
		admin.DELETE("/users/:userId", adminHandler.Delete)            //
		admin.GET("/users/:userId", adminHandler.FindByID)             //
		admin.GET("/users", adminHandler.ListUsers)                    //
		//admins
		admin.GET("/admins", adminHandler.ListAdmins) //
		//books
		admin.POST("/books", productHandler.Addproduct)                   //
		admin.PUT("/books/:bookId", productHandler.EditProductDet)        //
		admin.GET("/books", productHandler.ListProducts)                  //
		admin.DELETE("/books/:bookId", productHandler.DeleteProduct)      //
		admin.POST("books/:bookId/cover", productHandler.AddBookCover)    //
		admin.GET("/books/:bookId/covers", productHandler.ListBookCovers) //
		//orders
		admin.GET("/allorders", orderHandler.GetAllOrders)
		admin.GET("/topten", orderHandler.GetTopTen)
		//sales chart
		admin.GET("/chart", orderHandler.GetChart)
		//category
		admin.POST("/categories", categoryHandler.AddCategory)                  //
		admin.PUT("/categories/:categoryId", categoryHandler.UpdateCategory)    //
		admin.DELETE("/categories/:categoryId", categoryHandler.DeleteCategory) //
		//coupon
		engine.POST("/coupon", couponHandler.CreateNewCoupon)
		engine.DELETE("/coupon/:id", couponHandler.DeleteCoupon)
		engine.GET("/coupon/all", couponHandler.ListAllCoupon)

	}

	/*users.GET("/readbook", userHandler.ReadBook)
	users.GET("/premium",userHandler.GetPremium)
	users.POST("/premium",userHandler.MakePremium)*/

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
