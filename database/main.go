package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"./db_connection"
	"strconv"
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

func saveComments(c chan []Comments)  {
	var comments []Comments

	db := db_connection.SetConnection()
	defer db.Close()

	comments = <- c

	for i := range comments{

		res, err := db.Query(fmt.Sprintf("INSERT into `comments` (`post_id`, `email`, `name`, `body`) " +
			"VALUES ('%d', '%s', '%s', '%s')", comments[i].PostId,comments[i].Email, comments[i].Name, comments[i].Body))

		if err != nil{
			panic(err)
		}
		res.Close()
	}
}

func savePosts(posts []Post)  {
	db := db_connection.SetConnection()
	defer db.Close()

	for i := range posts {

		res, err := db.Query(fmt.Sprintf("INSERT IGNORE  into `posts` (`user_id`, `title`, `body`) " +
			"VALUES ('%d', '%s', '%s')", posts[i].UserId,posts[i].Title, posts[i].Body) +
			" SELECT * FROM (SELECT " +strconv.Itoa(posts[i].UserId)+",'"+posts[i].Title + "','"+ posts[i].Body+"') AS tmp WHERE NOT EXISTS ( " +
				" SELECT `user_id`, `title`, `body` FROM `posts` WHERE user_id = "+strconv.Itoa(posts[i].UserId)+" and title="+posts[i].Title+" and body="+posts[i].Body+") LIMIT1")

		if err != nil{
			panic(err)
		}
		res.Close()
	}

}