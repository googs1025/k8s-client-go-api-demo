package patch

import (
	"context"
	"fmt"
	"k8s-api-practice/initclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"testing"
)

type JSONPatch1 struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
type JSONPatchList1 []*JSONPatch1

func AddJsonPatch1(jps ...*JSONPatch1) JSONPatchList1 {
	list := make([]*JSONPatch1, len(jps))
	for index, jp := range jps {
		list[index] = jp
	}
	return list
}

func TestPatchPractice3(t *testing.T) {
	ctx := context.Background()
	// 动态客户端
	dynamicClient := initclient.ClientSet.DynamicClient
	patchPost := AddJsonPatch1(&JSONPatch1{
		Op:   "add",
		Path: "/spec/template/spec/containers/-",
		Value: map[string]interface{}{
			"name":  "redis",
			"image": "redis:5-alpine",
		},
	})

	gvr := schema.GroupVersionResource{
		Resource: "deployments",
		Version:  "v1",
		Group:    "apps",
	}

	b, _ := json.Marshal(patchPost)
	_, err := dynamicClient.Resource(gvr).Namespace("default").
		Patch(ctx, "patch-deployment", types.JSONPatchType, b, metav1.PatchOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("动态客户端patch成功")
}
