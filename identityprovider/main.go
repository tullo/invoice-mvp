package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/tullo/invoice-mvp/identityprovider/secret"
)

// User holds user info.
type User struct {
	name  string
	admin bool
}

// CustomClaims adds name and admin claims.
type CustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func token(name string, password string) string {
	var u User
	var ok bool
	if u, ok = login(name, password); !ok {
		return "Login failed"
	}

	cc := CustomClaims{u.name, u.admin, jwt.StandardClaims{
		Subject: "invoice.mvp",
	}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	jwt, err := t.SignedString([]byte(secret.Shared))
	if err != nil {
		panic(err)
	}
	return jwt
}

func login(u string, p string) (User, bool) {
	return User{name: "Go Invoicer", admin: true}, true
}

func main() {
	fmt.Printf("%v\n", token("go", "s3cr3t"))
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiaW52b2ljZS5tdnAifQ.j_NUeC0VmuxvrV-B1cVevUJPuBoxzXx2qbdg38otdh0
	/*
		header:
		{
			"alg": "HS256",
			"typ": "JWT"
		}

		payload:
		{
			"name": "Go Invoicer",
			"admin": true,
			"sub": "invoice.mvp"
		}
	*/
}
