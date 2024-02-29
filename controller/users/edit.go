package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type editRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func Edit(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 99)
	if !ok {
		return
	}

	req := editRequest{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	userID := session.UserID
	fail := false

	if req.Username == nil && req.Password == nil {
		fail = true
	}

	username := ""
	password := ""

	if req.Username != nil {
		username = *req.Username

		if username == "" {
			fail = true
		}
	}

	if req.Password != nil {
		password = *req.Password

		if password == "" {
			fail = true
		}
	}

	if fail {
		utils.BadRequest(c)
		return
	}

	if password != "" {
		password = utils.MD5Salt(password)
	}

	user := models.GetUserByID(userID)
	if user == nil {
		utils.BadRequest(c)
		return
	}

	if username != "" {
		user.Username = username
	}

	if password != "" {
		user.Password = password
	}

	models.EditUser(*user)

	models.AddLog("修改用户信息", fmt.Sprintf("%s已修改用户信息", user.Email), session.UserID)

	utils.Status200(c, nil, "")
}
