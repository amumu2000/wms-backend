package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminDeleteReq struct {
	UserID *[]int64 `json:"user_id"`
}

func AdminDelete(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminDeleteReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	userIDList := *req.UserID

	emails := make([]string, 0)

	for _, userID := range userIDList {
		user := models.GetUserByID(userID)

		if user == nil {
			utils.BadRequest(c)
			return
		}

		emails = append(emails, user.Email)
	}

	models.DeleteUser(userIDList)

	for _, email := range emails {
		models.AddLog("删除用户", fmt.Sprintf("已删除%s用户", email), session.UserID)
	}

	utils.Status200(c, nil, "")
}
