package main

import (
	"context"
	"fmt"
	"k8s-api-practice/initclient"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"log"
)

// Informer is a generic informer that can be used to watch multiple resources
type Informer[T runtime.Object] struct {
	Client   dynamic.Interface
	Group    string
	Version  string
	Resource string
	Handler  func(context.Context, T, watch.EventType) error
}

// NewInformer creates a new informer that watches the given resource
func NewInformer[T runtime.Object](client dynamic.Interface, group string, kind string, version string, handler func(context.Context, T, watch.EventType) error) *Informer[T] {
	return &Informer[T]{
		Client:   client,
		Group:    group,
		Version:  version,
		Resource: kind,
		Handler:  handler,
	}
}

// Watch watches the resource and returns a channel that receives the events
func (i *Informer[T]) Watch(ctx context.Context, options metav1.ListOptions) (chan T, error) {
	var t T
	res := schema.GroupVersionResource{
		Group:    i.Group,
		Version:  i.Version,
		Resource: i.Resource,
	}

	w, err := i.Client.Resource(res).Namespace("default").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	eventChan := make(chan T)
	go func() {
		defer close(eventChan)
		for {
			select {
			case event := <-w.ResultChan():
				ee := event.Object.(*unstructured.Unstructured)

				b, _ := ee.MarshalJSON()
				if err != nil {
					fmt.Println("MarshalJSON: ", err)
					return
				}

				err = json.Unmarshal(b, &t)
				if err != nil {
					fmt.Println("Unmarshal: ", err)
					return
				}

				eventType := watch.EventType(event.Type)
				if err := i.Handler(ctx, t, eventType); err != nil {
					fmt.Println("ssssss", err)
				}
				eventChan <- t
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventChan, nil
}

// Example usage:
func main() {
	client := initclient.ClientSet.DynamicClient

	informer := NewInformer[*v1.Pod](client, "", "pods", "v1", func(ctx context.Context, pod *v1.Pod, eventType watch.EventType) error {
		switch eventType {
		case watch.Added:
			fmt.Printf("Pod added: %+v\n", pod)
		case watch.Modified:
			fmt.Printf("Pod updated: %+v\n", pod)
		case watch.Deleted:
			fmt.Printf("Pod deleted: %+v\n", pod)
		}
		return nil
	})
	podsChan, err := informer.Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("mjmmmmmmmmmmmm")
	go func() {
		for {
			select {
			case e := <-podsChan:
				fmt.Println("mjmmmmmmmmmmmm")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("aaaaaaaa: ", e)

			}
		}
	}()

	informer1 := NewInformer[*appsv1.Deployment](client, "apps", "deployments", "v1", func(ctx context.Context, pod *appsv1.Deployment, eventType watch.EventType) error {
		switch eventType {
		case watch.Added:
			fmt.Printf("Pod added: %+v\n", pod)
		case watch.Modified:
			fmt.Printf("Pod updated: %+v\n", pod)
		case watch.Deleted:
			fmt.Printf("Pod deleted: %+v\n", pod)
		}
		return nil
	})
	depsChan, err := informer1.Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("mjmmmmmmmmmmmm")
	go func() {
		for {
			select {
			case e := <-depsChan:
				fmt.Println("mjmmmmmmmmmmmm")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("")
				fmt.Println("aaaaaaaa: ", e)

			}
		}
	}()

	select {}

}
