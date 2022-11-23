package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"testing"
)

func getClient() *kubernetes.Clientset{
	config := &rest.Config {
		Host:"http://1.14.120.233:8009",
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func TestGetConfigmap(t *testing.T) {

	cm, err := getClient().CoreV1().ConfigMaps("default").
		Get(context.Background(),"mycm",v1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range cm .Data {
		fmt.Printf("key=%s,value=%s\n",k,v)
	}
}
