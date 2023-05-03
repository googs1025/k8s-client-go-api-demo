package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// ConfigMapClient 是一个用于操作ConfigMap的客户端
type ConfigMapClient struct {
	dynamicClient dynamic.Interface
}

// NewConfigMapClient 创建一个ConfigMapClient对象
func NewConfigMapClient(dynamicClient dynamic.Interface) *ConfigMapClient {
	return &ConfigMapClient{
		dynamicClient: dynamicClient,
	}
}

// Create 创建一个ConfigMap
func (c *ConfigMapClient) Create(namespace, name string, data map[string]string) error {
	configMap := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"data": data,
		},
	}

	_, err := c.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}).
		Namespace(namespace).
		Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("ConfigMap %s created\n", name)
	return nil
}

// Get 获取指定的ConfigMap
func (c *ConfigMapClient) Get(namespace, name string) (*unstructured.Unstructured, error) {
	result, err := c.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	configMapObj := &unstructured.Unstructured{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(result.Object, configMapObj)
	if err != nil {
		return nil, err
	}

	return configMapObj, nil
}

// Update 更新指定的ConfigMap
func (c *ConfigMapClient) Update(namespace, name string, data map[string]string) error {
	result, err := c.Get(namespace, name)
	if err != nil {
		return err
	}

	result.Object["data"] = data

	_, err = c.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}).
		Namespace(namespace).
		Update(context.TODO(), result, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("ConfigMap %s updated\n", name)
	return nil
}

// Delete 删除指定的ConfigMap
func (c *ConfigMapClient) Delete(namespace, name string) error {
	err := c.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}).
		Namespace(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("ConfigMap %s deleted\n", name)
	return nil
}

func main() {
	// 获取当前用户的kubeconfig文件路径
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	// 使用kubeconfig文件创建一个rest.Config对象
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	configMapClient := NewConfigMapClient(dynamicClient)

	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	err = configMapClient.Create("default", "my-configmap", data)
	if err != nil {
		fmt.Printf("Failed to create ConfigMap: %v\n", err)
		return
	}

	configMap, err := configMapClient.Get("default", "my-configmap")
	if err != nil {
		fmt.Printf("Failed to get ConfigMap: %v\n", err)
		return
	}
	fmt.Printf("ConfigMap: %v\n", configMap)

	newData := map[string]string{
		"key1": "new-value1",
		"key2": "new-value2",
	}
	err = configMapClient.Update("default", "my-configmap", newData)
	if err != nil {
		fmt.Printf("Failed to update ConfigMap: %v\n", err)
		return
	}

	configMap, err = configMapClient.Get("default", "my-configmap")
	if err != nil {
		fmt.Printf("Failed to get ConfigMap: %v\n", err)
		return
	}
	fmt.Printf("Updated ConfigMap: %v\n", configMap)

	err = configMapClient.Delete("default", "my-configmap")
	if err != nil {
		fmt.Printf("Failed to delete ConfigMap: %v\n", err)
		return
	}
}
