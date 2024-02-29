package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

var (
	logEnabled = false
)

func InitLog(enabled bool) {
	logEnabled = enabled
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware(c *gin.Context) {
	if logEnabled {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		fmt.Println("Response body: " + blw.body.String())
	}

}

func GinRequestLogMiddleware(c *gin.Context) {
	if logEnabled {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		fmt.Println("Request body: " + string(body))
	}
}
