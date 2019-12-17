package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"register": "ok",
	})
}

func loginEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"authenticationToken": "someVal",
	})
}

func confEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"confKey": "confVal",
	})
}
