package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func getClient2() *kubernetes.Clientset{
	config := &rest.Config {
		Host:"http://1.14.120.233:8009",
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return c
}


func main() {

	// 使用informer机制监听
	fact:=informers.NewSharedInformerFactory(getClient2(), 0)

	cmInformer:=fact.Core().V1().ConfigMaps()
	cmInformer.Informer().AddEventHandler(&CmHandler{})

	fmt.Println("-------------开始监听configmaps----------")
	fact.Start(wait.NeverStop)
	select {}

}

// 回调
type CmHandler struct{}
func(c *CmHandler) OnAdd(obj interface{}){}
func(c *CmHandler) OnUpdate(oldObj, newObj interface{}){
	if newObj.(*v1.ConfigMap).Name == "mycm" {
		log.Println("mycm发生了变化")
	}
}
func(c *CmHandler)	OnDelete(obj interface{}){}