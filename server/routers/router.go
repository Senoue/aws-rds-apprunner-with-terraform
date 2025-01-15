package v1

import (
	"fmt"
	"os"

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

	// 環境変数からベースURLを取得（設定されていない場合はデフォルト値）
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // デフォルトのローカルURL
	}

	swaggerURL := fmt.Sprintf("%s/swagger/doc.json", baseURL)
	url := ginSwagger.URL(swaggerURL)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Listen and serve on defined address
	r.Run()
}
