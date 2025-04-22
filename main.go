package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huuloc2026/go-to-do.git/common"
	"github.com/huuloc2026/go-to-do.git/config"
	"github.com/huuloc2026/go-to-do.git/database"
	"github.com/huuloc2026/go-to-do.git/modules/items/model"
	ginItem "github.com/huuloc2026/go-to-do.git/modules/items/transport/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/v1")
	config.LoadConfig()
	db := database.ConnectDB()
	{
		items := v1.Group("/items")
		{
			items.POST("/", ginItem.CreateItem(db))
			items.GET("/", ListItem(db))
			items.GET("/:id", GetItems(db))
			items.PATCH("/:id", UpdateItems(db))
			items.DELETE("/:id", DeleteItems(db))
		}
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}

func GetItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.ToDoItemCreation
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("id = ?", id).First(&itemData).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData))
	}
}

func UpdateItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		var updated model.ToDoItemUpdate

		if err := c.ShouldBind(&updated); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("id = ?", id).Updates(&updated).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func DeleteItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		var updated model.ToDoItemUpdate

		if err := db.Table(model.TodoItem{}.TableName()).Where("id = ?", id).Delete(&updated).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func ListItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var (
			listItem []model.TodoItem
			paging   common.Paging
		)

		// Parse query params page và limit
		if err := c.ShouldBindQuery(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paging.Paging()

		offset := (paging.Page - 1) * paging.Limit

		var total int64
		if err := db.Table(model.TodoItem{}.TableName()).Count(&total).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		paging.Total = total

		if err := db.Table(model.TodoItem{}.TableName()).
			Order("id DESC").
			Limit(paging.Limit).
			Offset(offset).
			Find(&listItem).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse[[]model.TodoItem, common.Paging, any](listItem, paging, nil))
	}
}
