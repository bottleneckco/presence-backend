package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/bottleneckco/statuses-backend/db"
	"github.com/bottleneckco/statuses-backend/model"
	"github.com/bottleneckco/statuses-backend/web"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Overload(fmt.Sprintf("%s.env", runtime.GOOS), ".env")

	if gin.Mode() == gin.DebugMode {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// Migrate
	db.DB.AutoMigrate(&model.User{}, &model.Booking{}, &model.Group{}, &model.Room{}, &model.Status{})
	db.DB.Model(&model.Status{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.DB.Model(&model.Booking{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.DB.Table("user_groups").AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.DB.Table("user_groups").AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")

	web.StartServer()
}
