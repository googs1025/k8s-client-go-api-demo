package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)
// 从k8s.io/api@v0.25.4/apps/v1/register.go 拷贝并修改

// GroupName is the group name use in this package
const GroupName = "stable.example.com" // 需要修改为CRD的group
// 注册自己的自定义资源
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"} // 需要修改为CRD的version

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder initializes a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	// 添加 CronTab 与 CronTabList 这两个资源到 scheme
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CronTab{},
		&CronTabList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
