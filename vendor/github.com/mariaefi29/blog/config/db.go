package config

import (
	"os"

	"gopkg.in/mgo.v2"
)

var (
	//SMTPEmail contains email of google smtp server
	SMTPEmail string
	//SMTPPassword contains password of google smtp server
	SMTPPassword string
)

// DB instance of MongoDB
var DB *mgo.Database

// Posts are posts in a blog
var Posts *mgo.Collection

// Comments are comments to posts in a blog
var Comments *mgo.Collection

// Emails are subscription emails
var Emails *mgo.Collection

func init() {
	//smtp server credentials
	SMTPEmail = os.Getenv("SMTP_EMAIL")
	SMTPPassword = os.Getenv("SMTP_PASSWORD")
	// get a mongo sessions
	//DB_CONNECTION_STRING = mongodb://localhost/blog (env variable)
	s, err := mgo.Dial(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("blog_maria_efimenko")
	Posts = DB.C("posts")
	Comments = DB.C("comments")
	Emails = DB.C("emails")
	index := mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}
	Emails.EnsureIndex(index)
	// fmt.Println("You connected to your mongo database.")
}
