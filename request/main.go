package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Base URL
	var postsUri string = "https://jsonplaceholder.typicode.com/posts/"

	//Make Get Request
	response, err := http.Get(postsUri)

	//Check For Request Error
	if err != nil{
		panic(err)
		return
	}

	// Run After Function executed
	defer response.Body.Close()

	//Read Response Body
	body, error := ioutil.ReadAll(response.Body)

	//Check For Errors
	if error != nil{
		panic(error)
		return
	}

	//Body To String
	sb := string(body)

	// Print Result
	fmt.Printf(sb)
}
