package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"log"
	"time"
)

func (ApiTest) VerifyT(c *gin.Context) {
	jwksURL := "http://192.168.100.151:31924/jw"
	jwtToken := c.Query("tk")
	// from JWKS URL
	keySet, err := fetchJWKS(jwksURL)
	if err != nil {
		log.Fatalf("Failed to fetch JWKS: %v", err)
	}

	// parse JWT Token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// obtain JWT  kid (Key ID)
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// from JWKS obtain pub info
		key, found := keySet.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("key not found in JWKS")
		}

		//  jwk.Key transfer *rsa.PublicKey
		var rawKey interface{}
		if err := key.Raw(&rawKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %v", err)
		}
		return rawKey, nil
	})

	if err != nil {
		fmt.Printf("Failed to parse token: %v \r\n", err)
		return
	}

	if token.Valid {
		fmt.Println("Token is valid", token.Header, token.Claims)
		c.JSON(200, "token ok")
	} else {
		fmt.Println("Token is invalid", token)
		c.JSON(500, "Token is invalid")
	}
}

// fetchJWKS
func fetchJWKS(url string) (jwk.Set, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	set, err := jwk.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	return set, nil
}
