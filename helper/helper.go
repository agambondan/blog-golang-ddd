package helper

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func FailOnError(c *gin.Context, httpStatus int, err error) {
	if err != nil {
		c.JSON(httpStatus, gin.H{"message": err.Error()})
	}
}

