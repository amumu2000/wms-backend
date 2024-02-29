package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
)

func Self(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 99)
	if !ok {
		return
	}

	userID := session.UserID

	user := models.GetUserByID(userID)

	utils.Status200(c, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}, "")
}
