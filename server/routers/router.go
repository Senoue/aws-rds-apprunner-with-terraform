package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Senoue/aws-rds-apprunner-with-terraform/usecase"

	_ "github.com/Senoue/aws-rds-apprunner-with-terraform/docs" // Swagger API documentation
)

// StartService initializes and starts the Gin HTTP server.
func StartService(authUsecase *usecase.AuthUsecase) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.POST("/login", authUsecase.Login)
	r.POST("/register", authUsecase.Register)

	authGroup := r.Group("/v1")
	authGroup.Use(authUsecase.AuthMiddleware)
	authGroup.GET("/userInfo", authUsecase.UserInfo)

	// Swagger endpoint
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Listen and serve on defined address
	r.Run()
}
