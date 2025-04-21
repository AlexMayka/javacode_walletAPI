package routes

import (
	_ "JavaCode/docs"
	"JavaCode/internal/controllers"
	"JavaCode/internal/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin engine with routes, middleware, and Swagger.
//
// It registers API version groups, binds handlers to endpoints,
// and returns the fully configured *gin.Engine instance.
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	controller := controllers.Controller{DB: db}

	apiV1Group := router.Group("/api/v1")
	apiV1Group.Use(middleware.Logger())
	{
		apiV1Group.GET("wallets/:WALLET_UUID", controller.GetBalanceHandler)
		apiV1Group.POST("wallet", controller.WalletOperationHandler)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
