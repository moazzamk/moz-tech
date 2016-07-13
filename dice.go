package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)


func crawl() {

	url := "https://www.dice.com/jobs?q=php"

	req, err := http.NewRequest("POST", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("errorr")
	}

	defer resp.Body.Close()


	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response" , string(body))
}


func main() {
	crawl()
}


