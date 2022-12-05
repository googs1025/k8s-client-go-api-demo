## 手动创建CRD自定义对象
### 项目思路：
1. 本项目基于k8s内置的代码生成器生成CRD自定义对象
2. 并对CRD对象创建出Controller

### tag注解
pkg/apis/stable/v1beta1/doc.go
```bigquery
// +k8s:deepcopy-gen=package # 为包中任何类型生成深拷贝方法
// +groupName=stable.example.com # 指定group名称

package v1beta1
```

pkg/apis/stable/v1beta1/type.go
```bigquery
// +genclient # 指定要生成clientSet代码
// +genclient:noStatus # 指定不要生成Status 
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object # 指定实现runtime.Object基类 (实现接口对象)

// 根据 CRD 定义 CronTab 结构体
type CronTab struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CronTabSpec `json:"spec"`
}

// +k8s:deepcopy-gen=false # 不生成深拷贝方法

type CronTabSpec struct {
	CronSpec string `json:"cronSpec"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object # 指定实现runtime.Object基类 (实现接口对象)

// CronTab 资源列表
type CronTabList struct {
	metav1.TypeMeta `json:",inline"`

	// 标准的 list metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CronTab `json:"items"`
}
```
添加资源对象到Scheme注册表中
pkg/apis/stable/v1beta1/register.go
```bigquery
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
```