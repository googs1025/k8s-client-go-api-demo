package main

import (
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
)

type User struct {
	Id   	  int    `json:"-"`
	Name 	  string `json:"name"`
	Age  	  int 	`json:"age"`
	QQ        string `json:"qq"`
	LastApply string  `json:"-"`
}

func (u *User) Empty() bool {
	return u.Age == 0 && u.QQ == ""
}

func (u *User) ToJson() []byte {
	res, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	return res
}



func main() {

	user1 := &User{
		Name: "jiang",
		Age: 19,
		QQ: "aaaaa",
	}
	user2 := &User{
		Name: "jiang",
		Age: 20,
		QQ: "ddkfjal;df",
	}
	// 原对象
	user3 := &User{
		Name: "jiang",
	}

	// patch 会返回后者与前者的差别
	patch, _ := jsonpatch.CreateMergePatch(user1.ToJson(), user2.ToJson())
	fmt.Println(string(patch))

	// 把patch 放入原对象
	newuser, _ :=jsonpatch.MergePatch(user3.ToJson(), patch)
	fmt.Println(string(newuser))
}
