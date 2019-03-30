package main

import (
	"regexp"
	"github.com/gin-gonic/gin"
	"fmt"
)

var reg = regexp.MustCompile("https?:\\/\\/127\\.0\\.0\\.1\\:?.*")

func CORS(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		origin := c.Request.Header.Get("Origin")
		r := reg.Copy()
		if r.MatchString(origin) {
			headers := c.Request.Header.Get("Access-Control-Request-Headers")
			c.Writer.Header().Set("Access-Control-Allow-Origin", CORS_ORIGIN_WHITELIST)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
			c.Data(200, "text/plain", []byte{})
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			fmt.Println("dd")
		} else {
			c.Data(403, "text/plain", []byte{})
			c.Abort()
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Next()
			fmt.Println("dwwd")
		}

		return
	}
}