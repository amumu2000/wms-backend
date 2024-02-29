package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type loginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func Login(c *gin.Context) {
	req := loginRequest{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	email := ""
	password := ""
	fail := false

	//check email
	if req.Email == nil {
		fail = true
	}

	//check password
	if req.Password == nil {
		fail = true
	}

	if fail {
		utils.BadRequest(c)
		return
	}

	email = *req.Email
	password = *req.Password

	password = utils.MD5Salt(password)

	users := models.GetUsers("", password, email, -1)

	if len(users) == 0 {
		models.AddLog("用户登录", fmt.Sprintf("%s登录失败", email), -1)
		utils.Error400(c, nil, "用户名或密码错误。")
		return
	}

	user := users[0]

	//generate token
	session := utils.Session{
		UserID: user.ID,
		Role:   user.Role,
	}
	token, err := utils.GenerateToken(session)
	if err != nil {
		utils.InternalServerError(c)
		return
	}

	models.AddLog("用户登录", fmt.Sprintf("%s登录成功", email), user.ID)

	utils.Status200(c, gin.H{
		"token":   token,
		"user_id": user.ID,
		"role":    user.Role,
	}, "")
}
