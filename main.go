package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

type ToDoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status"`
}
type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Paging() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10 // default limit
	}
}
func (TodoItem) TableName() string {
	return "todoTables"
}
func (ToDoItemUpdate) TableName() string {
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

func GetItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData ToDoItemCreation
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("id = ?", id).First(&itemData).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": itemData,
		})
	}
}

func UpdateItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		var updated ToDoItemUpdate

		if err := c.ShouldBind(&updated); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("id = ?", id).Updates(&updated).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteItems(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		var updated ToDoItemUpdate

		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Delete(&updated).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func ListItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var (
			listItem []TodoItem
			paging   Paging
		)

		// Parse query params page và limit
		if err := c.ShouldBindQuery(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paging.Paging()

		offset := (paging.Page - 1) * paging.Limit

		var total int64
		if err := db.Table(TodoItem{}.TableName()).Count(&total).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		paging.Total = total

		if err := db.Table(TodoItem{}.TableName()).
			Order("id DESC").
			Limit(paging.Limit).
			Offset(offset).
			Find(&listItem).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   listItem,
			"paging": paging,
		})
	}
}
