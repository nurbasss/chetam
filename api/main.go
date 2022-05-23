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
	// POST /api/posts/ - adding a post with url or text
	r.HandleFunc("/api/posts", isAuthorized(createPost)).Methods("POST")
	// DELETE /api/post/{POST_ID} - delete post by id
	r.HandleFunc("/api/posts/{post_id}", isAuthorized(deletePostByID)).Methods("DELETE")
	// GET /api/posts/ - get all posts
	r.HandleFunc("/api/posts", getAllPosts).Methods("GET")
	// GET /api/post/{POST_ID}/upvote - upvote by post id
	r.HandleFunc("/api/posts/{post_id}/upvote", isAuthorized(upvotePost)).Methods("GET")
	// GET /api/post/{POST_ID}/downvote - downvote by post id
	r.HandleFunc("/api/posts/{post_id}/downvote", isAuthorized(downvotePost)).Methods("GET")
	// POST /api/post/{POST_ID} - add comment to a post by id
	r.HandleFunc("/api/posts/{post_id}", isAuthorized(addComment)).Methods("POST")
	// DELETE /api/post/{POST_ID}/{COMMENT_ID} - delete comment from a post by id
	r.HandleFunc("/api/posts/{post_id}/{comment_id}", isAuthorized(deleteComment)).Methods("DELETE")

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
