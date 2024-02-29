package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminAddReq struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Role     *int    `json:"role"`
}

func AdminAdd(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminAddReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	username := *req.Username
	password := *req.Password
	email := *req.Email
	role := *req.Role

	users := models.GetUsers("", "", email, -1)
	if len(users) > 0 {
		utils.Error400(c, nil, "邮箱已存在。")
		return
	}

	password = utils.MD5Salt(password)

	userID := models.AddUser(username, password, email, role)

	models.AddLog("添加用户", fmt.Sprintf("已添加%s用户", email), session.UserID)

	utils.Status200(c, gin.H{
		"user_id": userID,
	}, "")
}
