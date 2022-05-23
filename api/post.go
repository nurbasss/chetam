package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	mux "github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//Post type
type Post struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Category         string        `json:"category" bson:"category"`
	Type             string        `json:"type" bson:"type"`
	Title            string        `json:"title" bson:"title"`
	URL              string        `json:"url,omitempty" bson:"url"`
	Text             string        `json:"text,omitempty" bson:"text"`
	Author           Author        `json:"author" bson:"author"`
	Comments         []Comment     `json:"comments" bson:"comments"`
	Created          string        `json:"created" bson:"create"`
	Score            int           `json:"scope" bson:"scope"`
	Views            int           `json:"views" bson:"views"`
	UpvotePercentage int           `json:"upvotePercentage" bson:"upvotePercentage"`
	Votes            []Vote        `json:"votes" bson:"votes"`
}

//Author type
type Author struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
}
type Vote struct {
	User bson.ObjectId `json:"user" bson:"_id,omitempty"`
	Vote int           `json:"vote" bson:"vote"`
}

//Comment type
type Comment struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author  Author        `json:"author" bson:"author"`
	Created string        `json:"created" bson:"create"`
	Body    string        `json:"body" bson:"body"`
	Comment string        `json:"comment,omitempty"`
}

func createPost(wResponse http.ResponseWriter, r *http.Request) {
	var (
		body     []byte
		response []byte
		db       *mgo.Database
		token    *jwt.Token
		err      error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(wResponse, http.StatusBadRequest, "Chetam bad request type")
		return
	}

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	post := &Post{}
	if err = json.Unmarshal(body, post); err != nil {
		jsonError(wResponse, http.StatusBadRequest, "Chetam bad incorrect body")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = returnToken(r); err != nil {
		jsonError(wResponse, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(wResponse, http.StatusUnauthorized, err.Error())
		return
	}

	if status, err := postModel.create(post, author); status != http.StatusOK || err != nil {
		jsonError(wResponse, status, err.Error())
		return
	}

	if response, err = json.Marshal(post); err != nil {
		jsonError(wResponse, http.StatusInternalServerError, err.Error())
		return
	}

	wResponse.Header().Set("Content-Type", "application/json")
	wResponse.WriteHeader(http.StatusOK)
	wResponse.Write(response)
	wResponse.Write([]byte("\n\n"))
}

func (a *Author) fullFromJWTToken(token *jwt.Token) error {
	if a == nil {
		return errors.New("Chetam avtora netu")
	}

	var (
		user []byte
		err  error
	)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("chetam none claims")
	}

	if user, err = json.Marshal(claims["user"]); err != nil {
		return err
	}

	if err = json.Unmarshal(user, a); err != nil {
		return err
	}

	return nil
}

func deletePostByID(w http.ResponseWriter, r *http.Request) {
	var (
		db    *mgo.Database
		token *jwt.Token
		err   error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "chetam bad request type")
		return
	}

	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		jsonError(w, http.StatusBadRequest, "chetam takogo posta net")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	if token, err = returnToken(r); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	author := &Author{}
	if err := author.fullFromJWTToken(token); err != nil {
		jsonError(w, http.StatusUnauthorized, err.Error())
		return
	}

	post := &Post{ID: bson.ObjectId(postID)}
	if status, err := postModel.deleteByID(post, author); status != http.StatusOK || err != nil {
		jsonMessage(w, status, err.Error())
		return
	}

	jsonMessage(w, http.StatusOK, "chetam post deleted success")
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	var (
		response []byte
		db       *mgo.Database
		err      error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "chetam bad request type")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	postModel := PostModel{
		DB: db,
	}

	posts, status, err := postModel.getAll()
	if status != http.StatusOK || err != nil {
		jsonError(w, status, err.Error())
		return
	}

	if response, err = json.Marshal(posts); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	w.Write([]byte("\n\n"))
}
