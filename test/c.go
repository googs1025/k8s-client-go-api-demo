package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func getKubernetesClient() (*kubernetes.Clientset, error) {
	// 获取当前用户的kubeconfig文件路径
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// 使用kubeconfig文件创建一个rest.Config对象
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// 创建一个Kubernetes的Clientset对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createConfigMap(clientset *kubernetes.Clientset, namespace, name string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Data: data,
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("ConfigMap %s created\n", name)
	return nil
}

func getConfigMap(clientset *kubernetes.Clientset, namespace, name string) (*corev1.ConfigMap, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

func updateConfigMap(clientset *kubernetes.Clientset, namespace, name string, data map[string]string) error {
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	configMap.Data = data
	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("ConfigMap %s updated\n", name)
	return nil
}

func deleteConfigMap(clientset *kubernetes.Clientset, namespace, name string) error {

	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("ConfigMap %s deleted\n", name)
	return nil
}

func main() {
	// 获取Kubernetes客户端
	clientset, err := getKubernetesClient()
	if err != nil {
		fmt.Printf("Failed to get Kubernetes client: %v\n", err)
		return
	}

	// ConfigMap的命名空间和名称
	namespace := "default"
	name := "your-configmap"

	// 创建ConfigMap
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	err = createConfigMap(clientset, namespace, name, data)
	if err != nil {
		fmt.Printf("Failed to create ConfigMap: %v\n", err)
		return
	}

	// 获取ConfigMap
	configMap, err := getConfigMap(clientset, namespace, name)
	if err != nil {
		fmt.Printf("Failed to get ConfigMap: %v\n", err)
		return
	}
	fmt.Printf("ConfigMap: %v\n", configMap)

	// 更新ConfigMap
	newData := map[string]string{
		"key1": "new-value1",
		"key2": "new-value2",
	}
	err = updateConfigMap(clientset, namespace, name, newData)
	if err != nil {
		fmt.Printf("Failed to update ConfigMap: %v\n", err)
		return
	}

	// 获取更新后的ConfigMap
	configMap, err = getConfigMap(clientset, namespace, name)
	if err != nil {
		fmt.Printf("Failed to get ConfigMap: %v\n", err)
		return
	}

	// 删除ConfigMap
	err = deleteConfigMap(clientset, namespace, name)
	if err != nil {
		fmt.Printf("Failed to delete ConfigMap: %v\n", err)
		return
	}
}
