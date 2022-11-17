package field

import (
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"testing"
)

func TestFieldUse(t *testing.T) {

	// 模拟 外界需要输入的变量
	myField := fields.Set{
		"aaa": "aaa",
		"bbb": "bbb",
		"ccc": "ccc",
	}
	myField2 := fields.Set{
		"aaa": "aaa",
	}

	//sel := fields.SelectorFromSet(myField)
	sel2 := fields.SelectorFromSet(myField2) // 模拟pod上的
	// 是否匹配，只要有一个匹配上，就算匹配
	if sel2.Matches(myField) {
		fmt.Printf("Selector %v 匹配 field set %v\n", sel2, myField)
	} else {
		panic("Selector should have matched field set")
	}

	// key1=value1
	sel2 = fields.OneTermEqualSelector("aaa", "aaa")
	if sel2.Matches(myField) {
		fmt.Printf("Selector %v matched field set %v\n", sel2, myField)
	} else {
		panic("Selector should have matched field set")
	}

	// field1=value1,field2=value2，两个都匹配上
	sel2 = fields.AndSelectors(
		fields.OneTermEqualSelector("aaa", "aaa"),
		fields.OneTermEqualSelector("ccc", "ccc"),
	)
	if sel2.Matches(myField) {
		fmt.Printf("Selector %v matched field set %v\n", sel2, myField)
	} else {
		panic("Selector should have not matched field set")
	}

	// 直接解析+匹配
	sel, err := fields.ParseSelector("aaa==aaa")
	if err != nil {
		panic(err.Error())
	}
	if sel.Matches(myField) {
		fmt.Printf("Selector %v matched field set %v\n", sel, myField)
	} else {
		panic("Selector should have matched field set")
	}

}

/*
 字段选择器使用范例。
 eventList, _ := k8s.ClientSet.CoreV1().Events(namespace).List(ctxEvent, metav1.ListOptions{
	  FieldSelector: fields.Set{
	   "involvedObject.kind":      "Pod",
	   "involvedObject.name":      podName,
	   "involvedObject.namespace": podNamespace,
	  }.AsSelector().String(),
 })

*/
