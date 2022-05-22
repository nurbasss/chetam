package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//PostModel type
type PostModel struct {
	DB *mgo.Database
}

// Create new post
func (postModel *PostModel) create(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("chetam post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("chetam post not initialized")
	}

	if author == nil {
		return http.StatusInternalServerError, errors.New("chetam author not initialized")
	}

	post.ID = bson.NewObjectId()
	post.Score = 1
	post.Views = 0
	post.UpvotePercentage = 100

	now := time.Now()
	post.Created = fmt.Sprintf("%sT%sZ", now.Format("2006-01-02"), now.Format("03:04:05.000"))

	post.Votes = []Vote{{
		User: author.ID,
		Vote: 1},
	}

	post.Author = Author{
		ID:       author.ID,
		Username: author.Username,
	}

	post.Comments = []Comment{}

	if err := postModel.DB.C("posts").Insert(post); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (postModel *PostModel) deleteByID(post *Post, author *Author) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("chetam post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("chetam post not initialized")
	}

	if status, err := postModel.checkExistsByID(post); err != nil {
		return status, err
	}

	if err := postModel.DB.C("posts").Remove(
		bson.M{
			"_id":        bson.ObjectIdHex(string(post.ID)),
			"author._id": bson.ObjectId(author.ID),
		}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (postModel *PostModel) checkExistsByID(post *Post) (int, error) {
	if postModel == nil || postModel.DB == nil {
		return http.StatusInternalServerError, errors.New("chetam post model type not initialized")
	}

	if post == nil {
		return http.StatusInternalServerError, errors.New("chetam post not initialized")
	}

	if count, err := postModel.DB.C("posts").Find(
		bson.M{
			"_id": bson.ObjectIdHex(string(post.ID)),
		}).Count(); count == 0 || err != nil {
		if count == 0 {
			return http.StatusNotFound, errors.New("chetam post ne naiden")
		} else if count > 1 {
			return http.StatusInternalServerError, errors.New("chetam mnogo takih postov")
		}
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
