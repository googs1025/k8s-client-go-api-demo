package pkg

import (
	"context"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // 兼容 Windows 系统
}

// K8sClient 创建k8s clientset
func K8sClient() *kubernetes.Clientset {

	kubeConfig := filepath.Join(
		homeDir(), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		klog.Error("init kubeConfig error: ", err)
		return nil
	}

	// clientSet客户端
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error("create clientSet error: ", err)
		return nil
	}
	return clientSet
}

// GetPodsByDeployment 根据传入Deployment获取当前"正在"使用的pod
func GetPodsByDeployment(depName, ns string) []v1.Pod {
	clientSet := K8sClient()
	deployment, err := clientSet.AppsV1().Deployments(ns).Get(context.TODO(),
		depName, metav1.GetOptions{})
	if err != nil {
		klog.Error("create clientSet error: ", err)
		return nil
	}
	rsIdList := getRsIdsByDeployment(deployment, clientSet)
	podsList := make([]v1.Pod, 0)
	for _, rs := range rsIdList {
		pods := getPodsByReplicaSet(rs, clientSet, ns)
		podsList = append(podsList, pods...)
	}

	return podsList
}

// getPodsByReplicaSet 根据传入的ReplicaSet查询到需要的pod
func getPodsByReplicaSet(rs appv1.ReplicaSet, clientSet *kubernetes.Clientset, ns string) []v1.Pod {
	pods, err := clientSet.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Error("list pod error: ", err)
		return nil
	}

	ret := make([]v1.Pod, 0)
	for _, p := range pods.Items {
		// 找到 pod OwnerReferences uid相同的
		if p.OwnerReferences != nil && len(p.OwnerReferences) == 1 {
			if p.OwnerReferences[0].UID == rs.UID {
				ret = append(ret, p)
			}
		}
	}
	return ret

}

// getRsIdsByDeployment 根据传入的dep，获取到相关连的rs列表(滚更后的ReplicaSet就没用了)
func getRsIdsByDeployment(dep *appv1.Deployment, clientSet *kubernetes.Clientset) []appv1.ReplicaSet {
	// 需要使用match labels过滤
	rsList, err := clientSet.AppsV1().ReplicaSets(dep.Namespace).
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.Set(dep.Spec.Selector.MatchLabels).String(),
		})
	if err != nil {
		klog.Error("list ReplicaSets error: ", err)
		return nil
	}

	ret := make([]appv1.ReplicaSet, 0)
	for _, rs := range rsList.Items {
		ret = append(ret, rs)
	}
	return ret
}


// UpgradePodImage 原地升级pod镜像
func UpgradePodByImage(pod *v1.Pod, images ...string) {
	clientSet := K8sClient()
	patchList := make([]*patchOperation, 0)
	for k, image := range images {
		p := &patchOperation{
			Op: "replace",
			Path: fmt.Sprintf("/spec/containers/%v/image", k),
			Value: image,
		}
		patchList = append(patchList, p)

	}
	patchBytes, err := json.Marshal(patchList)
	if err != nil {
		klog.Error(err)
		return
	}

	jsonPatch, err := jsonpatch.DecodePatch(patchBytes)
	if err != nil {
		klog.Error("DecodePatch error: ", err)
		return
	}
	jsonPatchBytes, err := json.Marshal(jsonPatch)
	if err != nil {
		klog.Error("json Marshal error: ", err)
		return
	}
	_, err = clientSet.CoreV1().Pods(pod.Namespace).
		Patch(context.TODO(), pod.Name, types.JSONPatchType,
			jsonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// RestartPodByImage 原地重启pod的方式
func RestartPodByImage(pod *v1.Pod) {
	restartImage := pod.Spec.Containers[0].Image
	clientSet := K8sClient()

	// 改成任意一个镜像
	randomImage := "nginx:1.18-alpine"
	if restartImage == randomImage {
		randomImage = "nginx:1.19-alpine"
	}
	patch := fmt.Sprintf(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "%v"}]`, randomImage)
	patchBytes := []byte(patch)


	jsonPatch, err := jsonpatch.DecodePatch(patchBytes)
	if err != nil {
		klog.Error("DecodePatch error: ", err)
		return
	}
	jsonPatchBytes, err := json.Marshal(jsonPatch)
	if err != nil {
		klog.Error("json Marshal error: ", err)
		return
	}
	_, err = clientSet.CoreV1().Pods(pod.Namespace).
		Patch(context.TODO(), pod.Name, types.JSONPatchType,
			jsonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	// 延迟
	time.Sleep(time.Second * 30)

	// 再次使用patch换回原来image
	restartPatch := fmt.Sprintf(`[{"op": "replace", "path": "/spec/containers/0/image", "value": "%v"}]`, restartImage)
	restartPatchBytes := []byte(restartPatch)

	restartJsonPatch, err := jsonpatch.DecodePatch(restartPatchBytes)
	if err != nil {
		klog.Error("DecodePatch error: ", err)
		return
	}
	restartJsonPatchBytes, err := json.Marshal(restartJsonPatch)
	if err != nil {
		klog.Error("json Marshal error: ", err)
		return
	}
	_, err = clientSet.CoreV1().Pods(pod.Namespace).
		Patch(context.TODO(), pod.Name, types.JSONPatchType,
		restartJsonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}


