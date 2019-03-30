package main

import (
	"net/http"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
)
var CORS_ORIGIN_WHITELIST = "http://127.0.0.1:3000"
	// Here was the problem indeed and it has to be http://localhost:3000, not http://localhost:3000/

func getTwitterAPI(at, ats string) *anaconda.TwitterApi {
    anaconda.SetConsumerKey(consumerKey)
    anaconda.SetConsumerSecret(consumerSecret)
    return anaconda.NewTwitterApi(at, ats)
}

type StatusJSON struct {
	Status string `json:"status"`
}


func Post(c *gin.Context) {
	fmt.Println("post")
	var jsonIn StatusJSON
	ret := c.Bind(&jsonIn)
	fmt.Println(ret)

	session := sessions.Default(c)
	v := session.Get("request_token")
	if v == nil {
		fmt.Println("sessino nil")
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	
	rt := v.(string)
	fmt.Println(rt)
	v = session.Get("request_token_secret")
	if v == nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	rts := v.(string)
	fmt.Println(rts)
	api := getTwitterAPI(rt, rts)
	fmt.Println(jsonIn.Status)
	api.PostTweet(jsonIn.Status, nil)
	
	c.Header("Access-Control-Allow-Origin", CORS_ORIGIN_WHITELIST)
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, nil)
	return
}