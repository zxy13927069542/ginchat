package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/pion/ion-log"
	"time"
)

type MyClaims struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

type TokenResp struct {
	Token        string `json:"token"`
	AccessExpire string `json:"accessExpire"`
	RefreshAfter string `json:"refreshAfter"`
}

//	GenToken() 生成Token
//	todo: token过期时间和token密钥改用配置文件
func GenToken(name string, password string) *TokenResp {
	now := time.Now()
	claims := MyClaims{
		UserName: name,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(10800 * time.Second)),		//	3 hours
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := tokenStruct.SignedString([]byte("XXXXXXXXXXXXXX"))
	if err != nil {
		log.Errorf(">>Generate token failed! Err: [%v]", err)
		return nil
	}
	return &TokenResp{
		Token: signedString,
		AccessExpire: now.Add(10800 * time.Second).Format("2006-01-02 15:04:05"),
		RefreshAfter: now.Add(10800 / 2 * time.Second).Format("2006-01-02 15:04:05"),
	}
}

//	ParseToken()
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("XXXXXXXXXXXXXX"), nil
	})
	if err != nil {
		log.Errorf(">>Parse token failed! Err: [%v]", err)
		return nil, err
	} else if claims, ok := token.Claims.(*MyClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}
