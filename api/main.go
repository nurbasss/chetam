package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"

	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	rand.Seed(42)

	r := mux.NewRouter()

	// POST /api/register - registering a new user and getting a JWT token
	r.HandleFunc("/api/register", userRegister).Methods("POST")
	// POST /api/login - log in as an existing user and get a JWT token
	r.HandleFunc("/api/login", userLogin).Methods("POST")

	fmt.Println("starting server at :8080")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func jsonErrorRegister(w http.ResponseWriter, status int, user *User, msg string) {
	var errs = Errors{[]Error{
		{Location: "body",
			Param: "username",
			Value: user.Username,
			Msg:   msg},
	}}
	resp, _ := json.Marshal(errs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func jsonMessage(w http.ResponseWriter, status int, msg string) {
	var errs = map[string]interface{}{
		"message": msg,
	}
	resp, _ := json.Marshal(errs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}
