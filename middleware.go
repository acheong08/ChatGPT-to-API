package main

import (
	"os"

	gin "github.com/gin-gonic/gin"
)

var ADMIN_PASSWORD string

func init() {
	ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")
	if ADMIN_PASSWORD == "" {
		ADMIN_PASSWORD = "TotallySecurePassword"
	}
}

func admin_check(c *gin.Context) {
	password := c.Request.Header.Get("Authorization")
	if password != ADMIN_PASSWORD {
		c.String(401, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}
