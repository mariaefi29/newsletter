package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/mariaefi29/blog/config"

	"gopkg.in/mgo.v2/bson"
)

//Comment Struct
type Comment struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	PostID      bson.ObjectId `json:"post_id" bson:"post_id"`
	Content     string        `json:"content" bson:"content" schema:"message"`
	Author      string        `json:"author" bson:"author" schema:"username"`
	Email       string        `json:"email" bson:"email" schema:"email"`
	Website     string        `json:"website" bson:"website" schema:"website"`
	CreatedAt   string        `json:"time" bson:"time"`
	ApprovedFlg int           `json:"approved_flg" bson:"approved_flg"` //pending or approved. Pending by default.
}

//CreateComment puts a comment to a post into a database
func CreateComment(r *http.Request, idstr string) (Post, error) {
	// get form values
	comment := Comment{}
	post := Post{}

	post, err1 := OnePost(idstr)
	if err1 != nil {
		return post, errors.New("fail to find a post to comment: " + err1.Error())
	}
	err2 := r.ParseForm()
	if err2 != nil {
		return post, errors.New("fail to parse a comment form: " + err2.Error())
	}
	comment.ID = bson.NewObjectId()
	comment.PostID = bson.ObjectIdHex(idstr)
	currentTime := time.Now()
	comment.CreatedAt = currentTime.Format("02.01.2006 15:04:05")
	comment.ApprovedFlg = 0
	err3 := decoder.Decode(&comment, r.PostForm)
	if err3 != nil {
		return post, errors.New("fail to decode form into a struct: " + err3.Error())
	}

	// validate form values
	if comment.Email == "" || comment.Author == "" || comment.Content == "" {
		return post, errors.New("400 bad request: all fields must be complete")
	}

	// insert values to a database
	err4 := config.Comments.Insert(comment)
	if err4 != nil {
		return post, errors.New("500 internal server error: " + err4.Error())
	}
	//update a post
	post.Comments = append(post.Comments, comment)
	post.CommentCnt = 0
	for _, v := range post.Comments {
		if v.ApprovedFlg == 1 {
			post.CommentCnt++
		}
	}
	err5 := config.Posts.Update(bson.M{"_id": post.ID}, &post)
	if err5 != nil {
		return post, errors.New("500 internal server error: " + err5.Error())
	}

	return post, nil
}
