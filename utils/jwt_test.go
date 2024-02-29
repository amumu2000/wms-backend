package utils

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	InitJWT("testtest")

	session := Session{
		UserID: 100,
		Role:   1,
	}

	token, err := GenerateToken(session)
	if err != nil {
		t.Errorf("GenerateToken error: %s\n", err.Error())
	}

	fmt.Println(token)

	session2, err := ParseToken(token)
	if err != nil {
		t.Errorf("ParseToken error: %s\n", err.Error())
	}

	if session.UserID != session2.UserID || session.Role != session2.Role {
		t.Errorf("session unmatch.\n")
	}
}
