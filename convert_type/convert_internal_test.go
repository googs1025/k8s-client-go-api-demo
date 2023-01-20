package convert_type

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/kubernetes/pkg/apis/apps"

)

func TestConvertInternalType(t *testing.T) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1beta1.SchemeGroupVersion, &appsv1beta1.Deployment{})
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &appsv1.Deployment{})
	scheme.AddKnownTypes(apps.SchemeGroupVersion, &appsv1.Deployment{})

	metav1.AddToGroupVersion(scheme, appsv1beta1.SchemeGroupVersion)
	metav1.AddToGroupVersion(scheme, appsv1.SchemeGroupVersion)

	v1beta1Deployment := &appsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
			APIVersion: "apps/v1beta1",
		},
	}

	// v1 ----> internal版本
	objInternal, err := scheme.ConvertToVersion(v1beta1Deployment, apps.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("GVK: ", objInternal.GetObjectKind().GroupVersionKind().String())


	// internal版本  ----> v1 版本
	objV1, err := scheme.ConvertToVersion(objInternal, appsv1.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("GVK: ", objV1.GetObjectKind().GroupVersionKind().String())


	v1Deployment, ok := objV1.(*appsv1.Deployment)
	if !ok {
		panic("got wrong type")
	}

	fmt.Println("GVK: ", v1beta1Deployment.GetObjectKind().GroupVersionKind().String())

}
