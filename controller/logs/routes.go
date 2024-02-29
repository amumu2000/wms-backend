package logs

import (
	"github.com/gin-gonic/gin"
)

func Init(r *gin.RouterGroup) {
	logsGroup := r.Group("/v1/logs")

	logsGroup.POST("/list", List)
}
