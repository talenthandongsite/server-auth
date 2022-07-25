package jwt

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	ID                 string `json:"id,omitempty" bson:"_id,omitempty"`
	Username           string `json:"username,omitempty" bson:",omitempty"`
	AccessControl      string `json:"accessControl,omitempty" bson:",omitempty"`
	jwt.StandardClaims        // 표준 토큰 Claims
}

type Jwt struct {
	Secret   []byte
	Duration time.Duration
}

type JwtConfig struct {
	Secret   []byte
	Duration time.Duration
}

func Init(config JwtConfig) (jwt *Jwt, err error) {
	return &Jwt{
		Secret:   config.Secret,
		Duration: config.Duration,
	}, nil
}

func (j *Jwt) ForgeToken(id string, username string, accessControl string) (string, time.Time, error) {
	// 토큰 유효시간 설정
	expirationTime := time.Now().Add(j.Duration)

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
	signedAuthToken, err := atoken.SignedString([]byte(j.Secret))

	if err != nil {
		log.Println(err)
		return "", expirationTime, err
	}

	return "Bearer " + signedAuthToken, expirationTime, nil
}

func (j *Jwt) OpenToken(tokenStr string) (*JWTClaims, error) {

	tokenFrag := strings.Split(tokenStr, " ")
	if len(tokenFrag) != 2 {
		err := errors.New("token is corrupted")
		return nil, err
	}
	token := tokenFrag[1]

	var claims JWTClaims
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
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
	if time.Now().After(time.UnixMilli(claims.ExpiresAt)) {
		err := errors.New("token expired")
		return nil, err
	}

	return &claims, nil
}
