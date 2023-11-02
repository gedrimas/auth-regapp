package helpers

import (

	"github.com/gin-gonic/gin"
)

func IsAdminUserType(c *gin.Context) bool {
	userType := c.GetString("user_type")
	if userType == "ADMIN" {
		return true
	}
	return false
}

