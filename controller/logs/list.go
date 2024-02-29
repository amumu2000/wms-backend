package logs

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	if !utils.CheckToken(c, 0) {
		return
	}

	logs := models.GetLogs()
	data := make([]gin.H, 0)

	for _, log := range logs {
		data = append(data, gin.H{
			"id":      log.ID,
			"title":   log.Title,
			"message": log.Message,
		})
	}

	utils.Status200(c, gin.H{
		"logs": data,
	}, "")
}
