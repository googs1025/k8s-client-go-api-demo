package pod

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sort"
	"strings"
)

func ListPod(client *kubernetes.Clientset, namespace string, filter string) ([]string, map[string]*v1.Pod, error) {
	names := make([]string, 0)
	m := make(map[string]*v1.Pod, 0)
	ctx := context.Background()

	// 字段选取器 label选取器的用法
	labelSet := labels.SelectorFromSet(labels.Set(map[string]string{"app": "nginx"}))
	listOptions := metav1.ListOptions{
		LabelSelector: labelSet.String(),
		FieldSelector: "status.phase=Running", // fmt.Sprintf("spec.ports[0].nodePort=%s", port)
	}

	podList, err := client.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return names, m, err
	}
	for _, pod := range podList.Items {
		obj := pod
		name := pod.Name

		if filter == "" || strings.Contains(name, filter) && pod.DeletionTimestamp == nil {
			names = append(names, name)
			m[name] = &obj
		}

	}
	sort.Strings(names)
	return names, m, nil
}

func IsPodReady(pod *v1.Pod) bool {
	phase := pod.Status.Phase
	if phase != v1.PodRunning || pod.DeletionTimestamp != nil {
		return false
	}
	return
}

func IsPodCompleted() {

}

func IsPodReadyConditionTrue(status v1.PodStatus) {

}

func PodStatus() {

}

func GetPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_, condition := GetPodCondition(&status, v1.PodReady)
	return condition
}

func GetPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}
