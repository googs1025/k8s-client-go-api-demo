package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// Getter is a wrapper around a cache.Store that provides Get and List methods.
type Getter[T runtime.Object] interface {
	Get(name string) (T, bool)
	GetWithNamespace(name, namespace string) (T, bool)
	List() []T
}

type Indexer[T runtime.Object] struct {
	store cache.Store
}

func (g *Indexer[T]) Get(name string) (t T, exists bool) {
	obj, exists, err := g.store.GetByKey(name)
	if err != nil {
		return t, false
	}
	if !exists {
		return t, false
	}
	return obj.(T), true
}

func (g *Indexer[T]) GetWithNamespace(name, namespace string) (t T, exists bool) {
	return g.Get(namespace + "/" + name)
}

func (g *Indexer[T]) List() (list []T) {
	for _, obj := range g.store.List() {
		list = append(list, obj.(T))
	}
	return list
}
