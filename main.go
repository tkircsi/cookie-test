package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	addr   *string
	remote *string
)

var mySigningKey []byte

type UserData struct {
	Name    string `json:"name"`
	FulkID  string `json:"fulkid"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Agree   bool   `json:"agree"`
	jwt.StandardClaims
}

var userData = UserData{
	Name:    "Kiss Csaba",
	Age:     45,
	Address: "1222 Budapest, Szent János utca 7.",
	Agree:   true,
	FulkID:  "FLK-001122",
}

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	mySigningKey = []byte(secret)
}

func main() {
	addr = flag.String("addr", ":5000", "HTTP Server address")
	remote = flag.String("remote", "http://localhost:5002/redirpage", "Remote redirect URL")
	flag.Parse()

	http.HandleFunc("/", home)
	http.HandleFunc("/redirpage", redirpage)
	http.HandleFunc("/redirect", redirect)
	log.Printf("server listen on port %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Home"))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().Add(10 * time.Minute)
	token, err := GenerateJWT(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "UserData", Value: token, Path: "/", HttpOnly: true, Expires: expire, MaxAge: 90000}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, *remote, http.StatusTemporaryRedirect)
}

func redirpage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("UserData")
	if err != nil {
		log.Printf("cant find cookie :/\r\n")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userData, err := ParseJWT(cookie.Value)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "UserData: %v", userData)
}

func GenerateJWT(u UserData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u)

	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		err = fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(jwtFromHeader string) (*UserData, error) {

	token, err := jwt.ParseWithClaims(
		jwtFromHeader,
		&UserData{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	userData, ok := token.Claims.(*UserData)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	return userData, nil
}
