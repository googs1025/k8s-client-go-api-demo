package main

import (
	"context"
	"fmt"
	"k8s-api-practice/initclient"
	appv1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/coordination/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"log"
	"strings"
	"time"
)

// convertUnstructuredToResource 将 Unstructured 对象转换为 k8s 对象
func convertUnstructuredToResource[T runtime.Object](unstructuredObj *unstructured.Unstructured) (T, error) {
	var t T
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

type ResourceEventHandler[T runtime.Object] struct {
	AddFunc    func(obj T)
	UpdateFunc func(oldObj T, newObj T)
	DeleteFunc func(obj T)
}

func (e *ResourceEventHandler[T]) OnAdd(obj interface{}) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		rr, _ := convertUnstructuredToResource[T](o)
		e.AddFunc(rr)
	}
}

func (e *ResourceEventHandler[T]) OnUpdate(oldObj, newObj interface{}) {
	var t, tt *unstructured.Unstructured
	var ok bool
	if t, ok = oldObj.(*unstructured.Unstructured); !ok {
		return
	}
	if tt, ok = newObj.(*unstructured.Unstructured); !ok {
		return
	}
	oldT, err := convertUnstructuredToResource[T](t)
	if err != nil {
		return
	}
	newT, err := convertUnstructuredToResource[T](tt)
	if err != nil {
		return
	}
	e.UpdateFunc(oldT, newT)

}

func (e *ResourceEventHandler[T]) OnDelete(obj interface{}) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		rr, _ := convertUnstructuredToResource[T](o)
		e.DeleteFunc(rr)
	}
}

func parseGVR(gvr string) schema.GroupVersionResource {
	var group, version, resource string
	gvList := strings.Split(gvr, "/")
	if len(gvList) < 3 {
		group = ""
		version = gvList[0]
		resource = gvList[1]
	} else {
		if gvList[0] == "core" {
			gvList[0] = ""
		}
		group, version, resource = gvList[0], gvList[1], gvList[2]
	}
	return schema.GroupVersionResource{
		Group: group, Version: version, Resource: resource,
	}
}

func main() {
	// dynamic客户端
	client := initclient.ClientSet.DynamicClient

	factory := dynamicinformer.NewDynamicSharedInformerFactory(client, 5*time.Second)
	depDynamicInformer := factory.ForResource(parseGVR("apps/v1/deployments"))
	// eventHandler 回调
	depHandler := &ResourceEventHandler[*appv1.Deployment]{
		AddFunc: func(dep *appv1.Deployment) {
			fmt.Println("on add dep:", dep.Name)
		},
		UpdateFunc: func(old *appv1.Deployment, new *appv1.Deployment) {
			fmt.Println("on update dep:", new.Name)
		},
		DeleteFunc: func(dep *appv1.Deployment) {
			fmt.Println("on delete dep:", dep.Name)
		},
	}

	depDynamicInformer.Informer().AddEventHandler(depHandler)

	podDynamicInformer := factory.ForResource(parseGVR("core/v1/pods"))
	// eventHandler 回调
	podHandler := &ResourceEventHandler[*v1.Pod]{
		AddFunc: func(pod *v1.Pod) {
			fmt.Println("on add pod:", pod.Name)
		},
		UpdateFunc: func(old *v1.Pod, new *v1.Pod) {
			fmt.Println("on update pod:", new.Name)
		},
		DeleteFunc: func(pod *v1.Pod) {
			fmt.Println("on delete pod:", pod.Name)
		},
	}

	podDynamicInformer.Informer().AddEventHandler(podHandler)

	leaseDynamicInformer := factory.ForResource(parseGVR("coordination.k8s.io/v1/leases"))
	// eventHandler 回调

	leaseHandler := &ResourceEventHandler[*v12.Lease]{
		AddFunc: func(pod *v12.Lease) {
			fmt.Println("on add lease:", pod.Name)
		},
		UpdateFunc: func(old *v12.Lease, new *v12.Lease) {
			fmt.Println("on update lease:", new.Name)
		},
		DeleteFunc: func(pod *v12.Lease) {
			fmt.Println("on delete lease:", pod.Name)
		},
	}

	leaseDynamicInformer.Informer().AddEventHandler(leaseHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("------开始使用informer监听------------")
	factory.Start(ctx.Done())

	for gvr, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Fatal(fmt.Sprintf("Failed to sync cache for resource %v", gvr))
		}
	}

	select {}
}
