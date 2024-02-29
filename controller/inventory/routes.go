package inventory

import "github.com/gin-gonic/gin"

func Init(r *gin.RouterGroup) {
	inventoryGroup := r.Group("/v1/inventory")

	inventoryGroup.POST("/in", In)
	inventoryGroup.POST("/out", Out)
	inventoryGroup.POST("/list", List)
}
