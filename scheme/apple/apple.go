package apple

import (
	"k8s-api-practice/scheme/runtime"
	"k8s-api-practice/scheme/scheme"
	"strings"
)


type Apple struct {
	ApiVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind 	   string  		`json:"kind" yaml:"kind"`
	Metadata   				`json:"metadata" yaml:"metadata"`
	Spec 	   AppleSpec    `json:"spec" yaml:"spec"`
	Status     AppleStatus  `json:"status" yaml:"status"`

}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type AppleSpec struct {
	Size   	   string			`json:"size" yaml:"size"`
	Price  	   string			`json:"price" yaml:"price"`
	Place      string			`json:"place" yaml:"place"`
	Color      string			`json:"color" yaml:"color"`
}

type AppleStatus struct {
	//CreateTime time.Time
	Status     string
}

type AppleList struct {
	Item []*Apple
}

func (f *Apple) SetGroupVersionKind(kind runtime.GroupVersionKind) {
	f.Kind = kind.Kind
	if kind.Group == "" {
		f.ApiVersion = kind.Version
	} else {
		f.ApiVersion = kind.Group + "/" + kind.Version
	}
}

func (f *Apple) GroupVersionKind() runtime.GroupVersionKind {
	res := strings.Split(f.ApiVersion, "/")
	var s runtime.GroupVersionKind
	if len(res) < 2 {
		s = runtime.GroupVersionKind{
			Group: "",
			Version: res[0],
			Kind: f.Kind,
		}

	} else {
		s = runtime.GroupVersionKind{
			Group: res[0],
			Version: res[1],
			Kind: f.Kind,
		}
	}
	return s
}



func (f *Apple) GetObjectKind(g runtime.GroupVersionKind) (runtime.ObjectKind, error) {
	f.SetGroupVersionKind(g)
	return &Apple{}, nil
}

var SchemeGroupVersion = runtime.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Apple"}

var (
	schemeBuilder      = scheme.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
)

func addKnownTypes(scheme *scheme.Scheme) error {
	f := &Apple{
		ApiVersion: "apps/v1",
		Kind: "Apple",
	}
	scheme.AddKnownTypes(SchemeGroupVersion, f)
	return nil
}





