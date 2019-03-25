package main

import (

	"net/http"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
)

func getTwitterAPI(at, ats string) *anaconda.TwitterApi {
    anaconda.SetConsumerKey(consumerKey)
    anaconda.SetConsumerSecret(consumerSecret)
    return anaconda.NewTwitterApi(at, ats)
}

type StatusJSON struct {
	Status string `json:"status"`
}


func Post(c *gin.Context) {
	var jsonIn StatusJSON
	ret := c.Bind(&jsonIn)
	fmt.Println(ret)

	session := sessions.Default(c)
	v := session.Get("request_token")
	if v == nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	rt := v.(string)

	v = session.Get("request_token_secret")
	if v == nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	rts := v.(string)
	api := getTwitterAPI(rt, rts)
	api.PostTweet(jsonIn.Status, nil)

	c.JSON(http.StatusOK, nil)
	return
}