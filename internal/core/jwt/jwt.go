package jwt

import (
	"encoding/json"
	"errors"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	jwtpackage "github.com/golang-jwt/jwt/v4"

	"masmaint/config"
)


type JwtPayload struct {
	jwtpackage.StandardClaims
	CustomClaims
}


func GenerateJWT(claims CustomClaims) (string, error) {
	pl := generatePayload(claims)

	return encodeJWT(pl)
}


func generatePayload(claims CustomClaims) JwtPayload {
	var pl JwtPayload

	pl.CustomClaims = claims
	pl.IssuedAt =  time.Now().Unix()
	pl.ExpiresAt = time.Now().Add(time.Second * JWT_EXPIRES).Unix()

	return pl
}


func encodeJWT(payload JwtPayload) (string, error) {
	cf := config.GetConfig()
	token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, payload)
	return token.SignedString([]byte(cf.JwtSecretKey))
}


func decodeJWT(encoded string) (*jwtpackage.Token, error) {
	cf := config.GetConfig()
	token, err := jwtpackage.Parse(encoded, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(cf.JwtSecretKey), nil
	})

	return token, err
} 


func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtStr, _ := extractTokenFromCookie(c)
		pl, err := jwtAuth(jwtStr)

		if err != nil {
			c.Redirect(303, "/login")
			c.Abort()
			return
		}
		c.Set(CONTEXT_KEY_PAYLOAD, pl)
		c.Next()
	}
}


func JwtAuthApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//jwtStr, _ := extractTokenFromRequestHeader(c)
		jwtStr, _ := extractTokenFromCookie(c)
		pl, err := jwtAuth(jwtStr)

		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set(CONTEXT_KEY_PAYLOAD, pl)
		c.Next()
	}
}


func jwtAuth(jwtStr string) (JwtPayload, error) {
	token, err := decodeJWT(jwtStr)

	if err != nil {
		return JwtPayload{}, err
	}

	return getPayload(token)
}


func getPayload(token *jwtpackage.Token) (JwtPayload, error) {
	var pl JwtPayload

	jsonString, err := json.Marshal(token.Claims.(jwtpackage.MapClaims))

	if err == nil {
		err = json.Unmarshal(jsonString, &pl)
	}

	return pl, err
}


func extractTokenFromCookie (c *gin.Context) (string, error) {
	return c.Cookie(COOKIE_KEY_JWT)
} 


func extractTokenFromRequestHeader (c *gin.Context) (string, error) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("Hint: Authorization")
	}else if strings.Index(tokenString, "Bearer ") != 0 {
		return "", errors.New("Hint: Bearer")
	}
	return strings.TrimSpace(tokenString[7:]), nil
}
