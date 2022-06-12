package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/talenthandongsite/server-auth/internal/durable"
	"github.com/talenthandongsite/server-auth/internal/handler"
	"github.com/talenthandongsite/server-auth/internal/repo"

	"github.com/dgrijalva/jwt-go"
	"encoding/json"
	"time"
)

const PORT string = "8080"

func main() {
	ctx := context.Background()

	fmt.Println("Starting Talent Server")

	client, err := durable.GetClient(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()

	repository := repo.InitUserRepo(client)

	handler := handler.InitUserHandler(repository)
	
	mux.HandleFunc("/admin/user", handler.HandleCreateRead)
	mux.HandleFunc("/admin/user/", handler.HandleUpdateDelete)
	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/auth", Auth)

	app := http.FileServer(http.Dir("web"))
	assets := http.FileServer(http.Dir("assets"))

	mux.Handle("/app", http.StripPrefix("/app", app))
	mux.Handle("/app/", http.StripPrefix("/app/", app))
	mux.Handle("/assets", http.StripPrefix("/assets", assets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.ListenAndServe(":"+PORT, mux)
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

  // jwt키 생성
var jwtKey = []byte("my_secret_key")
  

  // Create a struct to read the username and password from the request body
type Credentials struct {	
	Username string `json:"username"`
	Password string `json:"password"`
  }
  
  // Create a struct that will be encoded to a JWT.
  // We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	  Username string `json:"username"`
	  jwt.StandardClaims
  }
  
  // Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	  var creds Credentials
	  // Get the JSON body and decode into credentials
	  err := json.NewDecoder(r.Body).Decode(&creds)
	  if err != nil {
		  // If the structure of the body is wrong, return an HTTP error
		  w.WriteHeader(http.StatusBadRequest)
		  return
	  }
  
	  // Get the expected password from our in memory map
	  expectedPassword, ok := users[creds.Username]
  
	  // 패스워드 다를경우 Unauth 응답
	  if !ok || expectedPassword != creds.Password {
		  w.WriteHeader(http.StatusUnauthorized)
		  return
	  }
  
	  // 토큰 유효시간 설정
	  expirationTime := time.Now().Add(5 * time.Minute)
	  // JWT claims 생성 username과 유효시간 포함
	  claims := &Claims{
		  Username: creds.Username,
		  StandardClaims: jwt.StandardClaims{
			  // In JWT, the expiry time is expressed as unix milliseconds
			  ExpiresAt: expirationTime.Unix(),
		  },
	  }
  
	  // 토큰 생성
	  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	  // JWT string 생성
	  tokenString, err := token.SignedString(jwtKey)
	  if err != nil {
		  //토큰 생성중 에러 발생한 경우 InternalServerError
		  w.WriteHeader(http.StatusInternalServerError)
		  return
	  }
	  //토큰을 통해서 쿠키 설정
	  http.SetCookie(w, &http.Cookie{
		  Name:    "token",
		  Value:   tokenString,
		  Expires: expirationTime,
	  })
  }

func Auth(w http.ResponseWriter, r *http.Request) {
	// 요청으로부터 쿠키의 토큰 가져오기
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 인증 성공한 경우 Welcome message 출력
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
