package dispatch

import "github.com/gin-gonic/gin"

func Init(r *gin.RouterGroup) {
	dispatchGroup := r.Group("/v1/dispatch")

	dispatchGroup.POST("/list", List)
	dispatchGroup.POST("/add", Add)
	dispatchGroup.POST("/delete", Delete)
	dispatchGroup.POST("/edit", Edit)
	dispatchGroup.POST("/complete", Complete)
	dispatchGroup.POST("/reject", Reject)
	dispatchGroup.POST("/confirm", Confirm)
}
