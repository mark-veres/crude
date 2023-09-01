package crude

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	DB               *gorm.DB
	CreateMiddleware gin.HandlersChain
	ReadMiddleware   gin.HandlersChain
	UpdateMiddleware gin.HandlersChain
	DeleteMiddleware gin.HandlersChain
}

func Register[T interface{}](c *Config, r *gin.RouterGroup, name string) {
	r.POST(name+"/new", genCreateHandlers[T](c)...)
	r.POST(name+"/update", genUpdateHandlers[T](c)...)
	r.GET(name+"/delete", genDeleteHandlers[T](c)...)
	r.GET(name+"/list", genListHandlers[T](c)...)
	r.GET(name+"/by/:property", genByPropertyHandlers[T](c)...)
	r.GET(name+"/where/:property/:operator", genWhereHandlers[T](c, name)...)
}

func genCreateHandlers[T interface{}](c *Config) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			var m T
			if err := ctx.BindJSON(&m); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
				return
			}

			if err := c.DB.Create(&m).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "database error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "successfully created",
			})
		},
	}
	return append(handlers, c.CreateMiddleware...)
}

func genListHandlers[T interface{}](c *Config) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			var result []T
			var m T

			err := c.DB.Model(&m).Find(&result).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "database error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"result": result})
		},
	}
	return append(handlers, c.ReadMiddleware...)
}

func genUpdateHandlers[T interface{}](c *Config) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			var m T
			if err := ctx.BindJSON(&m); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
				return
			}

			if err := c.DB.Save(&m).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "database error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "successfully updated",
			})
		},
	}
	return append(handlers, c.UpdateMiddleware...)
}

func genDeleteHandlers[T interface{}](c *Config) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			var m T
			id := ctx.Query("id")

			err := c.DB.Delete(&m, id).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "deletion error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "successfully deleted record",
			})
		},
	}
	return append(handlers, c.DeleteMiddleware...)
}

func genByPropertyHandlers[T interface{}](c *Config) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			prop := ctx.Param("property")
			value := ctx.Query("value")

			m := make(map[string]interface{})
			m[prop] = value
			var result []T

			err := c.DB.Where(m).Find(&result).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "database error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"result": result})
		},
	}
	return append(handlers, c.ReadMiddleware...)
}

func genWhereHandlers[T interface{}](c *Config, name string) gin.HandlersChain {
	handlers := gin.HandlersChain{
		func(ctx *gin.Context) {
			property := ctx.Param("property")
			operator := ctx.Param("operator")
			value := ctx.Query("value")
			from := ctx.Query("from")
			to := ctx.Query("to")
			pattern := ctx.Query("pattern")

			var err error
			var result []T

			tableName := c.DB.NamingStrategy.ColumnName(name, property)

			switch operator {
			case "=", ">", "<", ">=", "<=", "!=":
				query := fmt.Sprintf("%s %s ?", tableName, operator)
				err = c.DB.Find(&result, query, value).Error
			case "between":
				query := fmt.Sprintf("%s BETWEEN ? AND ?", tableName)
				err = c.DB.Where(query, from, to).Find(&result).Error
			case "like":
				query := fmt.Sprintf("%s LIKE ?", tableName)
				err = c.DB.Where(query, pattern).Find(&result).Error
			default:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "invalid operator",
				})
				return
			}

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "database error",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"result": result})
		},
	}
	return append(handlers, c.ReadMiddleware...)
}
