package patch

import (
	"context"
	_ "embed"
	"fmt"
	"k8s-api-practice/initClient"
	"k8s.io/apimachinery/pkg/util/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"log"
)

type JSONPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
type JSONPatchList []*JSONPatch

func AddJsonPatch(jps ...*JSONPatch) JSONPatchList {
	list := make([]*JSONPatch, len(jps))
	for index, jp := range jps{
		list[index]=jp
	}
	return list
}

func TestPatchPractice2(t *testing.T) {
	ctx := context.Background()
	client := initClient.ClientSet.Client

	var mgo, err = client.AppsV1().Deployments("default").
		Get(ctx, "patch-deployment", metav1.GetOptions{})
	if err!=nil{log.Fatalln(err)}
	// v1.Deployment{}

	patchPost := AddJsonPatch(&JSONPatch{
		Op: "remove",
		Path: "/spec/template/spec/containers/1", // 注意 0是容器的第一个 1是容器的第二个
		Value: map[string]interface{}{
			"name":"redis",
			"image":"redis:5-alpine",
		},
	})


	b, _ := json.Marshal(patchPost)
	fmt.Println(string(b))
	_, err = client.AppsV1().Deployments(mgo.Namespace).
		Patch(ctx,mgo.Name,types.JSONPatchType, b,
			metav1.PatchOptions{})

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("patch操作成功 JSONPatchType方式")






















}
