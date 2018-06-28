package main

import (
	"fmt"
	"log"

	"github.com/mariaefi29/blog/config"

	"github.com/mariaefi29/blog/models"
	gomail "gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2/bson"
)

func sendNewPostInfo(sub []models.Email, post models.Post) error {
	d := gomail.NewDialer("smtp.gmail.com", 587, config.SMTPEmail, config.SMTPPassword)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	m := gomail.NewMessage()
	greeting := "<h3>Уважаемый подписчик!<br></h3>"
	signature := "<p>С уважением,<br>Мария Ефименко</p>"
	content := fmt.Sprintf("<p>На сайте <a href=\"www.marialife.com\">www.marialife.com</a> опубликован новый пост под названием <a href=\"http://marialife.com/posts/show/%s\"><strong>%s</strong></a> в категории <a href=\"http://marialife.com/category/%s\">%s!</a></p>", post.IDstr, post.Name, post.CategoryEng, post.Category)
	for _, r := range sub {
		m.SetHeader("From", "no-reply@example.com")
		m.SetAddressHeader("To", r.EmailAddress, "Уважаемый подписчик!")
		m.SetHeader("Subject", "Новый пост в блоге Марии Ефименко")
		m.SetBody("text/html", fmt.Sprintf("%s%s%s", greeting, content, signature))

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r.EmailAddress, err)
			return err
		}
		m.Reset()
	}
	return nil
}

func main() {
	posts, err := models.AllPosts()
	if err != nil {
		log.Println(err)
		return
	}
	emails := []models.Email{}

	err1 := config.Emails.Find(bson.M{}).All(&emails)
	if err != nil {
		log.Println(err1)
		return
	}
	err2 := sendNewPostInfo(emails, posts[0])
	if err2 != nil {
		log.Println(err2)
		return
	}
}
