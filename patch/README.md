## Patch操作更新k8s集群资源
### kubectl 操作
#### 加入patch yaml
1. 创建一个名为 patch-file.yaml
```
spec:
  template:
    spec:
      containers:
      - name: patch-demo-ctr-2
        image: redis
```
2. 
```bigquery
kubectl patch deployment patch-demo --patch-file patch-file.yaml
```
### 注意patch类型 （根据官方文档）
1. 策略性合并类的patch
```bigquery

你在前面的练习中所做的 patch 称为 策略性合并 patch（Strategic Merge Patch）。
请注意，patch 没有替换 containers 列表。相反，它向列表中添加了一个新 Container。
换句话说， patch 中的列表与现有列表合并。当你在列表中使用策略性合并 patch 时，并不总是这样。 
在某些情况下，列表是替换的，而不是合并的。
对于策略性合并 patch，列表可以根据其 patch 策略进行替换或合并。 
patch 策略由 Kubernetes 源代码中字段标记中的 patchStrategy 键的值指定。 
```
总之：需要如何判断是patch操作是合并或是替换，需要由源码内字段决定。

2. 使用 JSON 合并 patch
```bigquery
策略性合并 patch 不同于 JSON 合并 patch。 使用 JSON 合并 patch，如果你想更新列表，你必须指定整个新列表。新的列表完全取代现有的列表。
```
```
add,replace,remove 
字符串格式 [{ "op": "replace", "path": "/xxx/xxx", "value": "xxx" }]
```

#### 参考文档：
````bigquery
https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
````