package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminEditReq struct {
	UserID   *int64  `json:"user_id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Role     *int    `json:"role"`
}

func AdminEdit(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminEditReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	userID := *req.UserID

	username, password, email, role := "", "", "", -1

	if req.Username != nil {
		username = *req.Username
	}

	if req.Password != nil {
		password = *req.Password
	}

	if req.Email != nil {
		email = *req.Email
	}

	if req.Role != nil {
		role = *req.Role
	}

	if userID == session.UserID && role >= 0 {
		utils.Error400(c, nil, "管理员无法修改自身的角色。")
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

	if email != "" {
		user.Email = email
	}

	if role >= 0 {
		user.Role = role
	}

	models.EditUser(*user)

	models.AddLog("编辑用户", fmt.Sprintf("已编辑%s用户", user.Email), session.UserID)

	utils.Status200(c, nil, "")
}
