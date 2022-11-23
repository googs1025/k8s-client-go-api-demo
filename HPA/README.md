
### Metrics-Server
k8s里。可以通过Metrics-Server服务采集节点和Pod的内存、磁盘、CPU和网络的使用率等。

Metrics API 只可以查询当前的度量数据，并不保存历史数据。

Metrics API URI 为 /apis/metrics.k8s.io/。


必须部署 metrics-server 才能使用该 API，metrics-server 通过调用 Kubelet Summary API 获取数据

#### 部署 metrics-server
```bigquery
# 需要vpn下载
1. 下载镜像
docker pull --platform=linux/arm64 bitnami/metrics-server:0.4.1
docker images
2. 打成tar包
docker save -o bitnami/metrics-server:0.4.1
3. 上传云服务器 scp -r xxxxxx
4. load镜像
docker load < aaa.tar
5. 执行
kubectl apply -f components.yaml
```

#### 执行命令
```bigquery
kubectl top node
kubectl autoscale deployment ngx1 --min=2 --max=5 --cpu-percent=20
```

### api
```bigquery
autoscaling/v1      #只支持通过cpu伸缩
autoscaling/v2beta1 #支持通过cpu、内存 和自定义数据来进行伸缩。
autoscaling/v2beta2 #beta1的进一步(一般使用)
```


#### 测试
可以使用deployment压测 测试HPA功能。