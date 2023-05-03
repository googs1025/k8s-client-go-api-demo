package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"
)

func main() {
	r := gin.New()
	r.POST("/authorize", func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			klog.Error(err)
			c.AbortWithStatusJSON(400, rsp(false, err.Error()))

		}
		obj := &unstructured.Unstructured{}
		err = json.Unmarshal(b, &obj)
		if err != nil {
			klog.Error(err)
			c.AbortWithStatusJSON(400, rsp(false, err.Error()))

		}

		klog.Info(string(b))
		c.JSON(200, rsp(true, ""))
	})
	r.RunTLS(":9090", "/Users/zhenyu.jiang/Desktop/debug_kubernetes/kubernetes/certs/apiserver.crt", "/Users/zhenyu.jiang/Desktop/debug_kubernetes/kubernetes/certs/apiserver.key")
}

const (
	AccessApiVersion = "authorization.k8s.io/v1beta1"
	AccessKind       = "SubjectAccessReview"
)

func rsp(allowed bool, reason string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(AccessApiVersion)
	obj.SetKind(AccessKind)
	obj.Object["status"] = map[string]interface{}{
		"allowed": allowed,
		reason:    reason,
	}
	return obj
}
