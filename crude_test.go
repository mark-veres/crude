package crude_test

import (
	"github.com/gin-gonic/gin"
	"github.com/mark-veres/crude"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Name    string
	Content string
}

var db *gorm.DB
var r *gin.Engine
var c crude.Config

func init() {
	db, _ = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	db.AutoMigrate(&Post{})

	r = gin.Default()
	api := r.Group("/api")

	c = crude.Config{
		DB: db,
	}

	crude.Register[Post](&c, api, "posts")

	r.Run(":8080")
}
