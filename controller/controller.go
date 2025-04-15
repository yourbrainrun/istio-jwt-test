package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"istio-jwt-test/utils"
	"net/http"
)

type ApiTest struct {
}

func (ApiTest) Public(c *gin.Context) {
	c.JSON(200, "ok public api")
}

func (ApiTest) Jw(c *gin.Context) {
	jwks := utils.GenerateJWKS()
	c.JSON(http.StatusOK, jwks)
	return
}

func (ApiTest) GenerateToken(c *gin.Context) {
	generateJWT, err := utils.GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": generateJWT,
	})
}

// SecureHandler jwt protected api
func (ApiTest) SecureHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a secure endpoint, sidecar check token pass",
	})
}

func (ApiTest) VerifyToken(c *gin.Context) {
	fmt.Println("pp:token :", c.Query("tk"))
	validateJWT, err := utils.ValidateJWT(c.Query("tk"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, validateJWT)
}
