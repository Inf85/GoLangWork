package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const baseUri string = "https://jsonplaceholder.typicode.com/posts/"

func main()  {
	for i := 1; i <=5; i++ {
		go makeRequest(i)
	}
	var input string
	fmt.Scanln(&input)
}

func makeRequest(page int)  {
	pageStr := strconv.Itoa(page)

	response, err := http.Get(baseUri + pageStr)

	//Check For Request Error
	if err != nil{
		panic(err)
		return
	}

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
