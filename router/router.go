package router

import (
	"github.com/gin-gonic/gin"
	"istio-jwt-test/controller"
)

func Api(r *gin.Engine) {
	r.GET("/jw", new(controller.ApiTest).Jw)
	r.GET("/jwks", new(controller.ApiTest).VerifyT)
	r.GET("/public", new(controller.ApiTest).Public)
	r.GET("/token", new(controller.ApiTest).GenerateToken)
	r.GET("/secure", new(controller.ApiTest).SecureHandler)
	r.GET("/vf", new(controller.ApiTest).VerifyToken)
}
