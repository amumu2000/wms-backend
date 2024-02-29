package warehouse

import "github.com/gin-gonic/gin"

func Init(r *gin.RouterGroup) {
	whGroup := r.Group("/v1/warehouse")

	whGroup.POST("/list", List)

	whAdminGroup := r.Group("/v1/admin/warehouse")

	whAdminGroup.POST("/add", AdminAdd)
	whAdminGroup.POST("/delete", AdminDelete)
	whAdminGroup.POST("/edit", AdminEdit)
}
