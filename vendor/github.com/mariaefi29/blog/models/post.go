package models

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/schema"
	"github.com/mariaefi29/blog/config"
	"gopkg.in/mgo.v2/bson"
)

var decoder = schema.NewDecoder()

//Post Struct
type Post struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	IDstr         string        `json:"idstr" bson:"idstr,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Category      string        `json:"category" bson:"category"`
	CategoryEng   string        `json:"categoryeng" bson:"categoryeng"`
	Date          string        `json:"date" bson:"date"`
	Images        []string      `json:"images" bson:"images"`
	Author        string        `json:"author" bson:"author"`
	Content       []string      `json:"content" bson:"content"`
	Likes         int           `json:"likes" bson:"likes"`
	Comments      []Comment     `json:"comments" bson:"comments"`
	CommentCnt    int           `json:"comments_cnt" bson:"comments_cnt"`
	IsPopular     int           `json:"popular" bson:"popular"`
	NextPostID    bson.ObjectId `json:"next_id" bson:"next_post_id,omitempty"`
	NextPostIDstr string        `json:"next_idstr" bson:"next_post_idstr,omitempty"`
	PrevPostID    bson.ObjectId `json:"prev_id" bson:"prev_post_id,omitempty"`
	PrevPostIDstr string        `json:"prev_idstr" bson:"prev_post_idstr,omitempty"`
}

func reverse(s []Post) []Post {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//AllPosts retrieves all posts
func AllPosts() ([]Post, error) {
	posts := []Post{}
	err := config.Posts.Find(bson.M{}).All(&posts)
	if err != nil {
		return nil, err
	}

	reverse(posts)

	return posts, nil
}

//OnePost retrieves one post by id
func OnePost(postIDstr string) (Post, error) {
	post := Post{}
	posts, err := AllPosts()
	if err != nil {
		return post, err
	}
	for i := range posts {
		if posts[i].IDstr == postIDstr {
			post = posts[i]
		}
	}
	return post, nil
}

//PostsByCategory retrieves posts by category
func PostsByCategory(categoryEng string) ([]Post, error) {
	posts := []Post{}

	err := config.Posts.Find(bson.M{"categoryeng": categoryEng}).All(&posts)
	if err != nil {
		return nil, err
	}

	reverse(posts)

	return posts, nil
}

//PostLike adds one like to a post
func PostLike(post Post) (int, error) {
	newLike := post.Likes + 1
	post.Likes++
	err := config.Posts.Update(bson.M{"_id": post.ID}, &post)
	if err != nil {
		return 0, err
	}
	return newLike, nil
}

//CreatePost writes a script of creating a post. Change the parameters inside the func and call the func once
func CreatePost() (Post, error) {

	post := Post{}

	post.ID = bson.NewObjectId()
	post.IDstr = post.ID.Hex()
	post.Name = "Откуда брать мотивацию"
	post.Category = "Английский язык"
	post.CategoryEng = "english"
	post.Date = "11 мая"
	post.Images = append(post.Images, "img/Motivation1.jpeg", "")
	post.Author = "Maria Efimenko"
	post.Content = append(post.Content, "Равным образом укрепление и развитие структуры обеспечивает широкому кругу (специалистов) участие в формировании систем массового участия. Товарищи! консультация с широким активом способствует подготовки и реализации модели развития.", "С другой стороны постоянное информационно-пропагандистское обеспечение нашей деятельности играет важную роль в формировании дальнейших направлений развития. Повседневная практика показывает, что постоянный количественный рост и сфера нашей активности обеспечивает широкому кругу (специалистов) участие в формировании систем массового участия. Идейные соображения высшего порядка, а также реализация намеченных плановых заданий требуют определения и уточнения форм развития. Разнообразный и богатый опыт постоянный количественный рост и сфера нашей активности играет важную роль в формировании дальнейших направлений развития. Таким образом рамки и место обучения кадров представляет собой интересный эксперимент проверки дальнейших направлений развития. Равным образом укрепление и развитие структуры влечет за собой процесс внедрения и модернизации соответствующий условий активизации.")
	post.Likes = 0
	post.CommentCnt = 0
	post.IsPopular = 0
	// post.NextPostID = ""
	// post.NextPostIDstr = ""
	post.PrevPostID = bson.ObjectIdHex("5a292ab35ca981492baab854")
	post.PrevPostIDstr = "5a292ab35ca981492baab854"
	fmt.Println(post)
	// insert values
	err := config.Posts.Insert(post)
	if err != nil {
		return post, errors.New("500 internal server error: " + err.Error())
	}
	return post, nil
}

//DeletePost deletes a post from a database
func DeletePost(postID string) error {
	err := config.Posts.Remove(bson.M{"_id": bson.ObjectIdHex(postID)})
	if err != nil {
		return errors.New("500 internal server error: " + err.Error())
	}
	return nil
}

//UpdatePost updates the post with requested parameters via a query.
// Query is expected to be a list of key=value settings separated by ampersands or semicolons.
// A setting without an equals sign is interpreted as a key set to an empty value.
func UpdatePost(postID string, query string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	post := Post{}
	// decoder := schema.NewDecoder()
	//parsing parametres, which need to be updated
	values, err := url.ParseQuery(query)
	if err != nil {
		return errors.New("Invalid input" + err.Error())
	}
	// find the post which needs to be updated
	err1 := config.Posts.Find(bson.M{"_id": bson.ObjectIdHex(postID)}).One(&post)
	if err1 != nil {
		return errors.New("500 internal database error," + err1.Error())
	}
	fmt.Println(post)
	fmt.Println(values)
	for k, v := range values {
		if k == "Content" || k == "Images" {
			for text != "q" { // break the loop if text == "q"
				fmt.Print("Enter your text: ")
				scanner.Scan()
				text = scanner.Text()
				if text != "q" {
					values.Add(k, text)
				}
			}
		}
		if k == "NextPostIDstr" {
			if v[0] != "" {
				post.NextPostID = bson.ObjectIdHex(v[0])
			} else {
				post.NextPostIDstr = ""
				post.NextPostID = ""
			}
		}
		if k == "PrevPostIDstr" {
			if v[0] != "" {
				post.PrevPostID = bson.ObjectIdHex(v[0])
			} else {
				post.PrevPostIDstr = ""
				post.PrevPostID = ""
			}
		}
		if k == "Comments" && v[0] == "delete" {
			_, err5 := config.Comments.RemoveAll(bson.M{"post_id": bson.ObjectIdHex(postID)})
			if err5 != nil {
				return errors.New(err5.Error())
			}
			post.Comments = []Comment{}
			post.CommentCnt = 0
			values.Del("Comments")
		}
	}
	err4 := decoder.Decode(&post, values)
	if err4 != nil {
		return errors.New(err4.Error())
	}
	fmt.Println(values)
	fmt.Println(post)
	err3 := config.Posts.Update(bson.M{"_id": post.ID}, post)
	if err3 != nil {
		return errors.New("500 internal database error," + err3.Error())
	}

	return nil
}
