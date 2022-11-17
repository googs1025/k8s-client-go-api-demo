## 字段选择器（Field selectors）
#### 允许你根据一个或多个资源字段的值 筛选 Kubernetes 资源
#### field selectors 使用"."点来访问字段。
```bigquery
metadata.name=my-service
metadata.namespace!=default
status.phase=Pending
```

### 命令行事例
```bigquery
# 单一字段
kubectl get pods --field-selector status.phase=Running
# 等价没有写
kubectl get pods --field-selector ""
# 支持 =、== 和 !=
kubectl get pods --field-selector=status.phase!=Running,spec.restartPolicy=Always
-- # 多种资源类型
kubectl get statefulsets,services --all-namespaces --field-selector metadata.namespace!=default
```

### client-go调用事例
```bigquery
eventList, _ := k8s.ClientSet.CoreV1().Events(namespace).List(ctxEvent, metav1.ListOptions{
	  FieldSelector: fields.Set{
	   "involvedObject.kind":      "Pod",
	   "involvedObject.name":      podName,
	   "involvedObject.namespace": podNamespace,
	  }.AsSelector().String(),
 })
```