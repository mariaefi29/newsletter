package models

import (
	"errors"
	"net/http"

	"github.com/mariaefi29/blog/config"
	"gopkg.in/mgo.v2/bson"
)

//Email Struct
type Email struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	EmailAddress string        `json:"email" bson:"email"`
}

//CreateEmail puts email address into a database
func CreateEmail(r *http.Request) (Email, error) {
	// get form values
	email := Email{}
	email.EmailAddress = r.FormValue("email")
	email.ID = bson.NewObjectId()

	// validate form values
	if email.EmailAddress == "" {
		return email, errors.New("400 bad request: all fields must be complete")
	}

	// insert values
	err1 := config.Emails.Insert(email)
	if err1 != nil {
		return email, errors.New("500 internal server error: " + err1.Error())
	}
	return email, nil
}
