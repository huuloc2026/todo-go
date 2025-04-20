package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huuloc2026/go-to-do.git/config"
	"github.com/huuloc2026/go-to-do.git/database"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id" gorm:"column:id;"`
	Title       string     `json:"title" gorm:"column:title;"`
	Description string     `json:"description" gorm:"column:description;"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

type ToDoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status"`
}

func (ToDoItemCreation) TableName() string {
	return "todoTables"
}

func main() {
	now := time.Now().UTC()
	item := TodoItem{
		Id:          1,
		Title:       "Task 1 Test",
		Description: "Content Description",
		Status:      "DOING",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	jsData, err := json.Marshal(item)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(jsData))
	r := gin.Default()
	v1 := r.Group("/v1")
	config.LoadConfig()
	db := database.ConnectDB()
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("/:id")
			items.PATCH("/:id")
			items.DELETE("/:id")
		}
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData ToDoItemCreation
		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&itemData).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": itemData.Id,
		})
	}
}
