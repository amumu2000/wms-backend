package goods

import "github.com/gin-gonic/gin"

func Init(r *gin.RouterGroup) {
	goodsGroup := r.Group("/v1/goods")

	goodsGroup.POST("/list", List)
	goodsGroup.POST("/add", Add)
	goodsGroup.POST("/delete", Delete)
	goodsGroup.POST("/edit", Edit)
}
