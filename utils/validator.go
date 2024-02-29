package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
)

func CheckToken(c *gin.Context, role int) bool {
	ret, _ := _CheckTokenWithSession(c, role)

	if !ret {
		TokenExpired(c)
	}
	return ret
}

func CheckTokenWithSession(c *gin.Context, role int) (bool, Session) {
	ret, session := _CheckTokenWithSession(c, role)

	if !ret {
		TokenExpired(c)
	}

	return ret, session
}

func _CheckTokenWithSession(c *gin.Context, role int) (bool, Session) {
	json := make(map[string]interface{})
	err := c.ShouldBindBodyWith(&json, binding.JSON)
	if err != nil {
		log.Printf("CheckToken error: %s\n", err.Error())
		return false, Session{}
	}

	tokenIntf, ok := json["token"]
	if !ok {
		log.Printf("CheckToken error: no token.\n")
		return false, Session{}
	}

	token, ok := tokenIntf.(string)
	if !ok {
		log.Printf("CheckToken error: token is not string.\n")
		return false, Session{}
	}

	session, err := ParseToken(token)
	if err != nil {
		{
			log.Printf("CheckToken error: parse token failed: %s\n", err.Error())
			return false, Session{}
		}
	}

	if session.Role > role {
		log.Printf("CheckToken error: role mismatch.\n")
		return false, session
	}

	return true, session
}
