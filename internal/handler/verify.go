package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	// 요청으로부터 쿠키의 토큰 가져오기
	token, ok := r.Header["Authorization"]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the JWT string from the cookie
	tknStr := token[0]
	tokenStr := strings.Split(tknStr, " ")
	log.Println(tokenStr[1])

	// Initialize a new instance of `Claims`
	claims, err := h.Jwt.OpenToken(tokenStr[1])
	log.Println(claims, err)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 인증 성공한 경우 Welcome message 출력
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
