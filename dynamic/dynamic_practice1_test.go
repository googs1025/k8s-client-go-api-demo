package dynamic

import (
	"context"
	"fmt"
	"k8s-api-practice/initclient"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"testing"
)

func TestDynamicClient1(t *testing.T) {

	client := initclient.ClientSet.DynamicClient

	namespace := "default"
	res := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}

	unList, err := client.Resource(res).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("list %s resource err: %s", res.Resource, err)
		return
	}

	// 遍历逻辑
	for _, item := range unList.Items {
		fmt.Println(item.GetName())
	}

	// 不定结构转为 GVK结构的方法。
	b, _ := unList.MarshalJSON()
	depList := v1.DeploymentList{}
	err = json.Unmarshal(b, &depList)
	if err != nil {
		return
	}

	for _, item := range depList.Items {
		fmt.Println(item.Kind)
	}

}
