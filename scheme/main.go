package main

import (
	"fmt"
	"k8s-api-practice/scheme/apple"
	"k8s-api-practice/scheme/product"
	"k8s-api-practice/scheme/runtime"
	"k8s-api-practice/scheme/scheme"
)

var sh = scheme.NewScheme()

// 初始化需要注册
var localSchemeBuilder = scheme.SchemeBuilder{
	product.AddToScheme,
	apple.AddToScheme,
}

var AddToScheme = localSchemeBuilder.AddScheme

func main() {
	err := AddToScheme(sh)
	if err != nil {
		return
	}
	var gvk = runtime.GroupVersionKind{Group: "food", Version: "v1", Kind: "Food"}
	res, err := sh.GetObjectKind(gvk)
	fmt.Println(res, err)

	var gvk1 = runtime.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Apple"}
	res1, err := sh.GetObjectKind(gvk1)
	fmt.Println(res1, err)

	aa, _ := res1.GetObjectKind(gvk1)

	fmt.Println(aa.GroupVersionKind())

}
