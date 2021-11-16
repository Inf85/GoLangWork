package main

import (
	"./db_connection"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const postsUri string = "https://jsonplaceholder.typicode.com/posts"
const commentsUri string = "https://jsonplaceholder.typicode.com/comments"
const userId string = "7"

type Post struct {
	UserId int
	Id int
	Title string
	Body string
}

type Comments struct {
	Id int
	PostId int
	Name string
	Email string
	Body string
}

func main() {
	var c chan []Comments = make(chan []Comments)
	posts := getPosts()
	savePosts(posts)
	for i := range posts{
		go getComments(posts[i].Id, c)
		go saveComments(c)
	}

	var input string
	fmt.Println("Press any key ...")
	fmt.Scanln(&input)
}

/**
Get Posts
 */
func getPosts() []Post {
	var post []Post
	response, err := http.Get(postsUri + "?userId=" + userId)

	if err != nil{
		panic(err)
		return nil
	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)
	//Check For Errors
	if error != nil{
		panic(error)
		return nil
	}
	errJson := json.Unmarshal(body, &post)
	if err != nil{
		panic(errJson)
		return nil
	}


	return post
}

/**
Get Comments
 */
func getComments(postId int, c chan []Comments)   {
	var comments []Comments
	var strPostId string = strconv.Itoa(postId)

	response, err := http.Get(commentsUri + "?postId=" + strPostId)

	if err != nil{
		panic(err)

	}
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)
	//Check For Errors
	if error != nil{
		panic(error)

	}
	errJson := json.Unmarshal(body, &comments)
	if err != nil{
		panic(errJson)

	}

	c <- comments
}

/**
Save Comment to Data Base
 */
func saveComments(c chan []Comments)  {
	var comments []Comments

	db := db_connection.SetConnection()
	defer db.Close()

	comments = <- c

	for i := range comments{
		stmt, err := db.Prepare(`INSERT into comments (id,post_id,email,name,body) VALUES (?,?,?,?,?)
		ON DUPLICATE KEY UPDATE post_id=?, email=?, name=?, body=?`)
		if err != nil{
			panic(err)
		}
		stmt.Exec(comments[i].Id,comments[i].PostId, comments[i].Email,comments[i].Name, comments[i].Body,comments[i].PostId, comments[i].Email,comments[i].Name, comments[i].Body)
	}
}

/**
Save Post to DataBase
 */
func savePosts(posts []Post)  {
	db := db_connection.SetConnection()
	defer db.Close()

	for i := range posts {
        var bodyText string =  strings.ReplaceAll(posts[i].Body,"\n","")
		fmt.Println(bodyText)

        stmt, err := db.Prepare(`INSERT into posts (id,user_id,title,body) VALUES (?,?,?,?)
		ON DUPLICATE KEY UPDATE user_id=?, title=?, body=?`)

		if err != nil{
			panic(err)
		}
		stmt.Exec(posts[i].Id,posts[i].UserId, posts[i].Title, bodyText,posts[i].UserId, posts[i].Title, bodyText)
	}

}