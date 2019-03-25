package main

import (
	"fmt"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginByTwitter(c *gin.Context) {
	oc := NewTWClient()
	rt, err := oc.RequestTemporaryCredentials(nil, callbackURL, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	session := sessions.Default(c)
	session.Set("request_token", rt.Token)
	session.Set("request_token_secret", rt.Secret)
	session.Save()

	url := oc.AuthorizationURL(rt, nil)

	c.Redirect(http.StatusMovedPermanently, url)
	return
}

func TwitterCallback(c *gin.Context) {
	tok := c.DefaultQuery("oauth_token", "")
	if tok == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	ov := c.DefaultQuery("oauth_verifier", "")
	if ov == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	session := sessions.Default(c)
	v := session.Get("request_token")
	if v == nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	rt := v.(string)
	if tok != rt {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	v = session.Get("request_token_secret")
	if v == nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	rts := v.(string)
	if rts == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	code, at, err := GetAccessToken(&oauth.Credentials{Token: rt, Secret: rts}, ov)
	if err != nil {
		fmt.Println(err)
		c.JSON(code, nil)
		return
	}

	session.Set("token", at.Token)
	session.Set("token_secret", at.Secret)
	session.Save()

	c.JSON(http.StatusOK, nil)
	return
}
