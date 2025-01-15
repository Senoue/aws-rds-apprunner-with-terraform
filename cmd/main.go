package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Senoue/aws-rds-apprunner-with-terraform/infrastructure"
	v1 "github.com/Senoue/aws-rds-apprunner-with-terraform/server/routers"
	"github.com/Senoue/aws-rds-apprunner-with-terraform/usecase"
	_ "github.com/go-sql-driver/mysql"
)

//	@title			Swagger Example API
//	@version		0.0.1
//	@description	This is a sample server.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	Cloudsmith inc.
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@BasePath		/v1
//	@schemes		http
//	@schemes		https

func main() {
	db := dbInit()
	defer db.Close()

	authUsecase := setupAuthComponents(db)
	v1.StartService(authUsecase)
}

func dbInit() *sql.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return sqldb
}

func setupAuthComponents(db *sql.DB) *usecase.AuthUsecase {
	authRepo := infrastructure.NewAuthRepository(db)
	return usecase.NewAuthUsecase(authRepo)
}
