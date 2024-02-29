package users

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type adminListReq struct {
	UserID   *int64  `json:"user_id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Role     *int    `json:"role"`
}

func AdminList(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 0)
	if !ok {
		return
	}

	req := adminListReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	userID, name, email, role := int64(-1), "", "", -1

	if req.UserID != nil {
		userID = *req.UserID
	}

	if req.Username != nil {
		name = *req.Username
	}

	if req.Email != nil {
		email = *req.Email
	}

	if req.Role != nil {
		role = *req.Role
	}

	users := models.FindUsers(userID, name, email, role)

	data := make([]gin.H, 0)
	for _, user := range users {
		if user.ID == session.UserID {
			continue
		}

		data = append(data, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"password": user.Password,
			"email":    user.Email,
			"role":     user.Role,
		})
	}

	utils.Status200(c, gin.H{
		"users": data,
	}, "")
}
