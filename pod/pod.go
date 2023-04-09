package pod

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	tools_watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sort"
	"strings"
	"time"
)


// ListPod list出所有pod，返回 pod名 对象map error
func ListPod(client *kubernetes.Clientset, namespace string, filter string) ([]string, map[string]*v1.Pod, error) {
	names := make([]string, 0)
	m := make(map[string]*v1.Pod, 0)
	ctx := context.Background()

	// 字段选取器 label选取器的用法
	labelSet := labels.SelectorFromSet(labels.Set(map[string]string{"app": "nginx"}))
	// list字段使用方法
	listOptions := metav1.ListOptions{
		LabelSelector: labelSet.String(),
		FieldSelector: "status.phase=Running", // fmt.Sprintf("spec.ports[0].nodePort=%s", port)
		Limit:         500,
	}

	// 从api server查询
	podList, err := client.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return names, m, err
	}
	for _, pod := range podList.Items {
		obj := pod
		name := pod.Name
		// 过滤
		if filter == "" || strings.Contains(name, filter) && pod.DeletionTimestamp == nil {
			names = append(names, name)
			m[name] = &obj
		}

	}
	sort.Strings(names)
	return names, m, nil
}

// IsPodReady 查看pod是否ready
func IsPodReady(pod *v1.Pod) bool {
	phase := pod.Status.Phase // Phase字段
	// 搭配metadata.finalizers使用，pod.DeletionTimestamp != nil表示没有被删除
	if phase != v1.PodRunning || pod.DeletionTimestamp != nil {
		return false
	}
	return IsPodReadyConditionTrue(pod.Status)
}

// IsPodCompleted 查看pod是否完成
func IsPodCompleted(pod *v1.Pod) bool {
	phase := pod.Status.Phase
	// 成功或失败条件就是pod已经完成并退出
	if phase == v1.PodSucceeded || phase == v1.PodFailed {
		return true
	}
	return false
}

func IsPodReadyConditionTrue(status v1.PodStatus) bool {
	condition := GetPodReadyCondition(status)
	return condition != nil && condition.Status == v1.ConditionTrue
}

// PodStatus 返回pod状态
func PodStatus(pod *v1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	phase := pod.Status.Phase
	if IsPodReady(pod) {
		return "Ready"
	}
	return string(phase)
}

// GetPodReadyCondition 返回pod ready状态
func GetPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_, condition := GetPodCondition(&status, v1.PodReady)
	return condition
}

// GetPodCondition 取得pod状态
func GetPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if status == nil {
		return -1, nil
	}
	// 注意需要遍例状态，这是一个列表
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

// 废弃，不用这个方法
func WaitForPodSelector(client *kubernetes.Clientset, namespace string, options metav1.ListOptions,
	timeout time.Duration) error {
	w, err := client.CoreV1().Pods(namespace).Watch(context.Background(), options)

	if err != nil {
		return err
	}
	defer w.Stop()

	condition := func(event watch.Event) (bool, error) {
		podObj := event.Object.(*v1.Pod)
		startTime := podObj.Status.StartTime
		complete := startTime != nil && !startTime.IsZero()
		if complete {
			//_, _ = yaml.Marshal(podObj)
			fmt.Printf("Pod %s is start: %s", podObj.Name, podObj.Status.Phase)
		}
		return complete, nil
	}

	//conditions := []ConditionFunc{
	// func(event watch.Event) (bool, error) {
	//  return event.Type == watch.Added, nil
	// },
	// func(event watch.Event) (bool, error) {
	//  return event.Type == watch.Modified, nil
	// },
	//}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	_, err = tools_watch.UntilWithoutRetry(ctx, w, condition)

	if err == wait.ErrWaitTimeout {
		return fmt.Errorf("pod %s never became ready", options.String())
	}
	return nil
}

// GetPodNames list出所有pod，并保存names
func GetPodNames(client kubernetes.Interface, ns string, filter string) ([]string, error) {
	names := make([]string, 0)
	list, err := client.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return names, fmt.Errorf("Failed to load Pods %s", err)
	}
	for _, d := range list.Items {
		name := d.Name
		if filter == "" || strings.Contains(name, filter) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

// GetPodRestarts 获取pod container重试次数
func GetPodRestarts(pod *v1.Pod) int32 {
	var restarts int32
	statuses := pod.Status.ContainerStatuses
	if len(statuses) == 0 {
		return restarts
	}
	for _, status := range statuses {
		restarts += status.RestartCount
	}
	return restarts
}