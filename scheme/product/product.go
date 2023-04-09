package product

import (
	"k8s-api-practice/scheme/runtime"
	"k8s-api-practice/scheme/scheme"
	"strings"
)


type Food struct {
	ApiVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind 	   string  		`json:"kind" yaml:"kind"`
	Metadata   				`json:"metadata" yaml:"metadata"`
	Spec 	   FoodSpec    `json:"spec" yaml:"spec"`
	Status     FoodStatus  `json:"status" yaml:"status"`

}

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type FoodSpec struct {
	Size   	   string			`json:"size" yaml:"size"`
	Price  	   string			`json:"price" yaml:"price"`
	Place      string			`json:"place" yaml:"place"`
	Color      string			`json:"color" yaml:"color"`
}

type FoodStatus struct {
	//CreateTime time.Time
	Status     string
}

type FoodList struct {
	Item []*Food
}

func (f *Food) SetGroupVersionKind(kind runtime.GroupVersionKind) {
	f.Kind = kind.Kind
	if kind.Group == "" {
		f.ApiVersion = kind.Version
	} else {
		f.ApiVersion = kind.Group + "/" + kind.Version
	}

}

func (f *Food) GroupVersionKind() runtime.GroupVersionKind {
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



func (f *Food) GetObjectKind(g runtime.GroupVersionKind) (runtime.ObjectKind, error) {
	f.SetGroupVersionKind(g)
	return &Food{}, nil
}

var SchemeGroupVersion = runtime.GroupVersionKind{Group: "food", Version: "v1", Kind: "Food"}

var (
	schemeBuilder      = scheme.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
)

func addKnownTypes(scheme *scheme.Scheme) error {
	f := &Food{
		ApiVersion: "v1",
		Kind: "Food",
	}
	scheme.AddKnownTypes(SchemeGroupVersion, f)
	return nil
}
