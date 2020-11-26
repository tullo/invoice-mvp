package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go/v4"
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
		Subject: "f8c39a31-9ced-4761-8a33-b9c628a67510", // User ID
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
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiR28gSW52b2ljZXIiLCJhZG1pbiI6dHJ1ZSwic3ViIjoiZjhjMzlhMzEtOWNlZC00NzYxLThhMzMtYjljNjI4YTY3NTEwIn0.WI6cRXYnYqUAV6qqNtf4B8PdGMgKuHqENQP5N_iCZL8
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
			"sub": "f8c39a31-9ced-4761-8a33-b9c628a67510"
		}
	*/
}
