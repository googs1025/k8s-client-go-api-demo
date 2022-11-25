package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

/*
	请求，测试HPA用
 */

func request() {
	// 请求。
	req, _ := http.NewRequest("GET",
		"http://42.193.17.123:31880/user/visit?userid=aaaaa", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func main() {
	for {
		request()
		time.Sleep(time.Second * 1)
	}
}
