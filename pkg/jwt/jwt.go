package jwt

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/talenthandongsite/server-auth/pkg/variable"
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

func Init(ctx context.Context) (jwt *Jwt, err error) {
	sstr := variable.GetEnv(ctx, variable.TOKEN_SECRET)
	secret := []byte(sstr)

	dstr := variable.GetEnv(ctx, variable.TOKEN_DURATION)
	dint, err := strconv.ParseInt(dstr, 10, 64)
	if err != nil {
		return nil, err
	}

	duration := time.Millisecond * time.Duration(dint)

	return &Jwt{
		Secret:   secret,
		Duration: duration,
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

	return signedAuthToken, expirationTime, nil
}

func (j *Jwt) OpenToken(tokenStr string) (*JWTClaims, error) {
	var claims JWTClaims
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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

	return &claims, nil
}
