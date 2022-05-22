package main

import (
	"errors"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//userRepository type
type UserModel struct {
	DB *mgo.Database
}
//for registration new user
func (um *UserModel) register(user *User) (int, error) {
	if um == nil || um.DB == nil {
		return http.StatusInternalServerError, errors.New("chetam user model not initialized")
	}

	if user == nil {
		return http.StatusInternalServerError, errors.New("chetam user not initialized")
	}
	//checks if username is unique
	if count, err := um.DB.C("users").Find(
		bson.M{"username": user.Username}).Count(); count > 0 || err != nil {
		if count > 0 {
			return http.StatusInternalServerError, errors.New("chetam takoi user uzhe exists")
		}
		return http.StatusInternalServerError, err
	}

	user.ID = bson.NewObjectId()
	user.Password, _ = hashPassword(user.Password)

	if err := um.DB.C("users").Insert(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (um *UserModel) login(user *User) (int, error) {
	if um == nil || um.DB == nil {
		return http.StatusInternalServerError, errors.New("chetam user model type not initialized")
	}

	if user == nil {
		return http.StatusInternalServerError, errors.New("chetam user not initialized")
	}
	//checks if user with that username exists
	if count, err := um.DB.C("users").Find(bson.M{
		"username": user.Username,
	}).Count(); count == 0 || err != nil {
		if count == 0 {
			return http.StatusNotFound, errors.New("chetam user ne naiden")
		} else if count > 1 {
			return http.StatusInternalServerError, errors.New("chetam mnogo takih userov")
		}
		return http.StatusInternalServerError, err
	}

	var u User
	if err := um.DB.C("users").Find(
		bson.M{"username": user.Username}).One(&u); err != nil {
		return http.StatusInternalServerError, err
	}

	if match := checkPasswordHash(user.Password, u.Password); !match {
		return http.StatusInternalServerError, errors.New("chetam invalid password")
	}

	*user = *&u

	return http.StatusOK, nil
}
