package ginItem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huuloc2026/go-to-do.git/common"
	"github.com/huuloc2026/go-to-do.git/modules/items/biz"
	"github.com/huuloc2026/go-to-do.git/modules/items/model"
	"github.com/huuloc2026/go-to-do.git/modules/items/storage"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.ToDoItemCreation
		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		store := storage.NewStore(db)
		bussiness := biz.CreateItemStorage(store)
		if err := bussiness.CreateItem(c.Request.Context(), &itemData); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}
