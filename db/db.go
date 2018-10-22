package db

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PQ driver
	"github.com/joho/godotenv"
)

var (
	// DB database object
	DB *gorm.DB
)

func init() {
	godotenv.Overload(fmt.Sprintf("%s.env", runtime.GOOS), ".env")

	var (
		dbUser    = os.Getenv("DB_USER")
		dbPass    = os.Getenv("DB_PASS")
		dbHost    = os.Getenv("DB_HOST")
		dbName    = os.Getenv("DB_NAME")
		dbURL     = os.Getenv("DATABASE_URL")
		dbSSLMode = os.Getenv("DB_SSL")
	)

	if len(dbUser) == 0 {
		dbUser = os.Getenv("USER")
	}
	var err error
	if len(dbURL) > 0 {
		DB, err = gorm.Open("postgres", dbURL)
	} else {
		DB, err = gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s sslmode=%s password=%s", dbUser, dbName, dbHost, dbSSLMode, dbPass))
	}
	DB.LogMode(gin.Mode() == gin.DebugMode)

	if err != nil {
		log.Panic(err)
	}

}
