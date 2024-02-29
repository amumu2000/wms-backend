package users

import "github.com/gin-gonic/gin"

func Init(r *gin.RouterGroup) {
	usersGroup := r.Group("/v1/users")

	usersGroup.POST("/login", Login)
	usersGroup.POST("/self", Self)
	usersGroup.POST("/edit", Edit)

	usersAdminGroup := r.Group("/v1/admin/users")
	usersAdminGroup.POST("/list", AdminList)
	usersAdminGroup.POST("/delete", AdminDelete)
	usersAdminGroup.POST("/add", AdminAdd)
	usersAdminGroup.POST("/edit", AdminEdit)
}
