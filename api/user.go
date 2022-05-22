package main

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	bson "gopkg.in/mgo.v2/bson"
)

type User struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func (u *User) getJWTToken() (string, error) {
	if u == nil {
		return "", errors.New("user type not initialized")
	}

	var (
		tokenString string
		err error
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{} {
			"username": u.Username,
			"id": u.ID,
		},
		"iat": time.Now().Add(time.Hour * 0).Unix(), // iat = issued at
		"exp": time.Now().Add(time.Hour * 24).Unix(), // exp = expiration time
	})

	if tokenString, err = token.SignedString(tokenSecret); err != nil {
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(wResponse http.ResponseWriter, request *http.Request) {
		if request.Header[tokenName] != nil {
			if _, err := returnToken(request); err != nil {
				jsonError(wResponse, http.StatusUnauthorized, err.Error())
				return
			}

			endpoint(wResponse, request)
		} else {
			jsonError(wResponse, http.StatusUnauthorized, "chetam ne avtorizovan")
			return
		}
	})
}