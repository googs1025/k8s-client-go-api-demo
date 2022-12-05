package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// 根据 CRD 定义 CronTab 结构体
type CronTab struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CronTabSpec `json:"spec"`
}

// +k8s:deepcopy-gen=false

type CronTabSpec struct {
	CronSpec string `json:"cronSpec"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronTab 资源列表
type CronTabList struct {
	metav1.TypeMeta `json:",inline"`

	// 标准的 list metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CronTab `json:"items"`
}