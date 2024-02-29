package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	signingKey string
)

type Session struct {
	UserID int64
	Role   int
}

func InitJWT(key string) {
	signingKey = key
}

func GenerateToken(session Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": session.UserID,
		"role":    session.Role,
		"exp":     time.Now().AddDate(0, 0, 7).Unix(),
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(tokenString string) (Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(signingKey), nil
	})

	if err != nil {
		return Session{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDIntf, roleIntf := claims["user_id"], claims["role"]
		userID, role := int64(-1), -1

		if userIDTmp, ok := userIDIntf.(float64); !ok {
			return Session{}, fmt.Errorf("userIDIntf.(int) not ok")
		} else {
			userID = int64(userIDTmp)
		}

		if roleTmp, ok := roleIntf.(float64); !ok {
			return Session{}, fmt.Errorf("roleIntf.(int) not ok")
		} else {
			role = int(roleTmp)
		}

		session := Session{
			UserID: userID,
			Role:   role,
		}

		return session, nil
	} else {
		return Session{}, TokenNotValidError
	}
}

type tokenNotValidError struct {
	Message string
}

func (t tokenNotValidError) Error() string {
	return t.Message
}

var (
	TokenNotValidError = tokenNotValidError{
		Message: "token not valid",
	}
)
