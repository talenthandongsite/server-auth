package jwtservice

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")
var tokenDuration time.Duration = time.Hour * 72

type JWTClaims struct {
	ID                 string `json:"id,omitempty" bson:"_id,omitempty"`
	Username           string `json:"username,omitempty" bson:",omitempty"`
	AccessControl      string `json:"accessControl,omitempty" bson:",omitempty"`
	jwt.StandardClaims        // 표준 토큰 Claims
}

type JwtService struct {
	Secret string
}

func Init(secret string) *JwtService {
	return &JwtService{
		Secret: secret,
	}
}

func (j *JwtService) ForgeToken(id string, username string, accessControl string, tokenDuration time.Duration) (string, time.Time, error) {
	// 토큰 유효시간 설정
	expirationTime := time.Now().Add(tokenDuration)

	// JWT claims 생성 username과 유효시간 포함
	claims := JWTClaims{
		ID:            id,
		Username:      username,
		AccessControl: accessControl,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.UnixMilli(),
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedAuthToken, err := atoken.SignedString(jwtKey)

	if err != nil {
		log.Println(err)
		return "", expirationTime, err
	}

	return signedAuthToken, expirationTime, nil
}

func (j *JwtService) OpenToken(tokenStr string) (*JWTClaims, error) {
	var claims JWTClaims
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return &claims, nil
}
