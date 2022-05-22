package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

func userLogin(wResponse http.ResponseWriter, request *http.Request) {
	var (
		body []byte
		response []byte
		db *mgo.Database
		tokenString string
		err error
	)

	if request.Header.Get("Content-Type") != "application/json" {
		jsonError(wResponse, http.StatusBadRequest, "chetam bad request type for login")
		return
	}

	if body, err = ioutil.ReadAll(request.Body); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error()+ "\n chetam body request kaida?")
		return
	}
	defer request.Body.Close()

	user := &User{}
	if err = json.Unmarshal(body, user); err != nil {
		jsonError(wResponse, http.StatusBadRequest, "chetam can not convert body to user")
		return
	}

	if user.Username == "" || user.Password == "" {
		jsonError(wResponse, http.StatusUnauthorized, "chetam login or password is empty")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	um := UserModel{
		DB: db,
	}

	if status, err := um.login(user); status != http.StatusOK {
		jsonError(wResponse, status, err.Error())
		return
	}

	if tokenString, err = user.getJWTToken(); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}

	if response, err = json.Marshal(map[string]interface{}{
		"token": tokenString,
	}); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}

	wResponse.Header().Set("Content-Type", "application/json")
	wResponse.WriteHeader(http.StatusCreated)
	wResponse.Write(response)
	wResponse.Write([]byte("\n\n"))
}