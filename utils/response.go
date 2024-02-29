package utils

import "github.com/gin-gonic/gin"

func Status200(c *gin.Context, data any, message string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	if message == "" {
		message = "成功"
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

func Error400(c *gin.Context, data any, message string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	if message == "" {
		message = "非法请求"
	}

	c.JSON(200, gin.H{
		"code":    400,
		"message": message,
		"data":    data,
	})
}

func Error401(c *gin.Context, data any, message string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	if message == "" {
		message = "登录失效，请重新登录。"
	}

	c.Header("content-type", "application/json; charset=utf-8")
	c.JSON(200, gin.H{
		"code":    401,
		"message": message,
		"data":    data,
	})
}

func Error500(c *gin.Context, data any, message string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	if message == "" {
		message = "内部服务器错误。"
	}

	c.Header("content-type", "application/json; charset=utf-8")
	c.JSON(200, gin.H{
		"code":    500,
		"message": message,
		"data":    data,
	})
}

func TokenExpired(c *gin.Context) {
	Error401(c, nil, "登录失效，请重新登录。")
}

func BadRequest(c *gin.Context) {
	Error400(c, nil, "非法请求")
}

func InternalServerError(c *gin.Context) {
	Error500(c, nil, "内部服务器错误")
}
