package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var RedirectUrl = "http://127.0.0.1:3000"


func main() {

	store := sessions.NewCookieStore([]byte("secret"))

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://127.0.0.1:3000", "http://127.0.0.1:8080"}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("session", store))
	r.Use(cors.New(config))

	r.GET("/login/twitter/auth", LoginByTwitter)
	r.GET("/login/twitter/auth/callback", TwitterCallback)
	r.POST("/twitter/post", Post)

	r.Run("127.0.0.1:8080")
}
