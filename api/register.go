package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

func userRegister(wResponse http.ResponseWriter, request *http.Request) {
	var (
		body        []byte
		response    []byte
		db          *mgo.Database
		tokenString string
		err         error
	)

	if request.Header.Get("Content-Type") != "application/json" {
		jsonError(wResponse, http.StatusBadRequest, "chetam request type incorrect")
		return
	}

	if body, err = ioutil.ReadAll(request.Body); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}
	defer request.Body.Close()

	user := &User{}
	if err = json.Unmarshal(body, user); err != nil {
		jsonError(wResponse, http.StatusBadRequest, "chetam can not convert body")
		return
	}

	if user.Username == "" || user.Password == "" {
		jsonError(wResponse, http.StatusUnauthorized, "chetam login ili parol ne ukazal")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	userModel := UserModel{
		DB: db,
	}

	if status, err := userModel.register(user); status != http.StatusOK {
		jsonErrorRegister(wResponse, http.StatusUnprocessableEntity, user, err.Error())
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
