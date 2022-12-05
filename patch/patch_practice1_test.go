package patch

import (
	"context"
	"fmt"
	"k8s-api-practice/initClient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"testing"
)

func TestPatchPractice1(t *testing.T) {
	ctx := context.Background()
	client := initClient.ClientSet.Client //获取 clientset
	var patchDeployment, err = client.AppsV1().Deployments("default").
		Get(ctx, "patch-deployment", metav1.GetOptions{})
	if err != nil{
		log.Fatal(err)
	}


	// 范例一：修改副本数
	patchPost := map[string]interface{}{
		"spec":map[string]interface{}{
			"replicas":1,
		},
	}

	// 范例二：修改容器镜像，注意需要从spec后全部写出来，不然会报错。
	//patchPost := map[string]interface{}{
	//	"spec":map[string]interface{}{
	//		"template":map[string]interface{}{
	//			"spec":map[string]interface{}{
	//				"containers":[]map[string]interface{}{
	//					{
	//						"name":"redis",
	//						"image":"redis:5-alpine",
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

	// 删除操作
	//patchPost := map[string]interface{}{
	//	"spec":map[string]interface{}{
	//		"template":map[string]interface{}{
	//			"spec":map[string]interface{}{
	//				"containers":[]map[string]interface{}{
	//					{
	//						"name":"redis",
	//						"$patch":"delete", // 加入这个
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

	b, _ := json.Marshal(patchPost)
	// patch操错
	_, err = client.AppsV1().Deployments(patchDeployment.Namespace).
		Patch(ctx, patchDeployment.Name, types.StrategicMergePatchType,b,metav1.PatchOptions{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("patch操作成功 StrategicMergePatchType方式")
}
