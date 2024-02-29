package utils

import (
	"fmt"
	"testing"
)

func TestMD5Salt(t *testing.T) {
	salt = "my_crypto_salt"

	fmt.Println(MD5Salt("sample_md5"))
}
