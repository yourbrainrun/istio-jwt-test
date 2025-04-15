package main

import (
	"github.com/gin-gonic/gin"
	"istio-jwt-test/router"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	// define router
	router.Api(r)

	r.Run(":8888")
}
