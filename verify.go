package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func verify() {
	r, err := http.PostForm("https://rrdoc.ru/rest-auth/login/", url.Values{"username": {"test1"}, "password": {"q1w2e3r4"}})
	if err != nil {
	panic(err)
	}
	body_byte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
  	fmt.Println("response ", string(body_byte))
}

func main() {
  	verify()
}
