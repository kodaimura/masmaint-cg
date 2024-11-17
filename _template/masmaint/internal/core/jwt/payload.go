package jwt

import (
	"time"
	
	jwtpackage "github.com/golang-jwt/jwt/v4"
)


type Payload struct {
	jwtpackage.StandardClaims
	CustomClaims
}

type CustomClaims struct {
	AccountId int
	AccountName string
	/* 独自のフィールドを追加可能 */
}


func NewPayload(claims CustomClaims) Payload {
	var pl Payload

	pl.CustomClaims = claims
	pl.IssuedAt =  time.Now().Unix()
	pl.ExpiresAt = time.Now().Add(time.Second * JWT_EXPIRES).Unix()

	return pl
}