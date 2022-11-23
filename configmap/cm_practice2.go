package main

import (
	"context"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
)

var api_server string
var token string
func init() {
	api_server = fmt.Sprintf("https://%s:%s",
		os.Getenv("KUBERNETES_SERVICE_HOST"),os.Getenv("KUBERNETES_PORT_443_TCP_PORT"))
	f, err := os.Open("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Fatal(err)
	}
	b, _ := ioutil.ReadAll(f)
	token = string(b)
}
func getClient1() *kubernetes.Clientset{
	config := &rest.Config {
		//Host:"http://124.70.204.12:8009",
		Host:api_server,
		BearerToken:token,
		TLSClientConfig:rest.TLSClientConfig{
			CAFile:"/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
func main() {
	cm, err := getClient1().CoreV1().ConfigMaps("default").
		Get(context.Background(),"mycm",v1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range cm.Data {
		fmt.Printf("key=%s,value=%s\n",k,v)
	}
	select {}
}
