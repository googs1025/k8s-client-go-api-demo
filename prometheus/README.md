## prometheus部署与应用练习练习


### 部署prometheus容器
```bigquer
#基于docker的方式进行部署
# 端口9090，运行后可通过浏览器 <云服务器ip>:9090进行访问
# prometheus.yml 配置文件需要先创建，如下。
docker run -d --name pm -p 9090:9090 -v /home/prometheus_config/config:/config     prom/prometheus:v2.30.0 --web.enable-lifecycle --config.file=/config/prometheus.yml 
```
prometheus.yml配置文件内容如下。
```bigquery
global:
  scrape_interval:     20s # 默认抓取间隔, 20秒向目标抓数据

```
### 部署 kube-state-metrics
官方提供的k8s内部各个组件的状态指标，通过监听 API Server 生成有关资源对象的状态指标， 如 Deployment、Node、Pod等。
默认有 /metrics 暴露http服务，供prometheus抓取.
```bigquery
pod为例：
kube_pod_info
kube_pod_status_ready
kube_pod_status_scheduled
kube_pod_container_status_terminated_reason
```
#### 部署方式
```bigquery
# 进入prometheus/state-metrics
kubectl apply -f .

# 部署好后可放问32280端口
```
#### 重新配置prometheus.yml文件
```bigquery
scrape_configs:
  - job_name: 'prometheus-state-metrics'
    static_configs:
    - targets: ['<云服务器ip(有部署kube-state-metrics的集群)>:32280']
# reload组件
curl -X POST http://<ip>:9090/-/reload  
```

### 部署node_exporter
监控CPU、内存、磁盘、I/O等信息，可以使用node-exporter。
#### 部署方式
```bigquery
# 进入prometheus/node_exporter
kubectl apply -f install.yaml
```

#### 重新配置prometheus.yml文件
```bigquery
global:
  scrape_interval: 20s

scrape_configs:
  - job_name: 'prometheus-state-metrics'
    static_configs:
    - targets: ['<ip>:32280']
  - job_name: 'node-exporter'
    static_configs:
    - targets: ['<ip>:9100','<ip>:9100']
    
# 重启服务 curl -X POST http://<ip>:9090/-/reload
```


### prometheus服务自动发现

Prometheus添加被监控端 两种方式：

• 静态配置：手动配置

• 服务发现：动态发现需要监控的 实例
其中服务发现支持来源有consul_sd_configs、file_sd_configs、kubernetes_sd_configs

Promethues通过k8s API集成目前主要支持5种服务发现模式，分别是：Node、Service、Pod、Endpoints、Ingress。

#### 操作步骤
1. 创建一个serviceaccount、clusterrole和clusterrolebinding 用于外部访问

2. 拷贝k8s集群的ca证书
```bigque
cd /etc/kubernetes/pki
cp ca.cert /home/prometheus_config/confi
```
3. 拷贝serviceaccount对应的token内容
```bigquery
# 取得token内容
kubectl -n kube-system describe secret \
$(kubectl -n kube-system describe sa  myprometheus  |grep  'Mountable secrets'| cut -f 2- -d ":" | tr -d " ") |grep -E '^token' | cut -f2 -d':' | tr -d '\t'

```
保存到 prometheus对应的服务器(/home/prometheus_config/config/)中(本项目使用同集群)

4. 增加prometheus.yml文件
```bigquery
- job_name: 'k8s-node'
    metrics_path: /metrics
    kubernetes_sd_configs:
      - api_server: https://<master ip>:6443/
        role: node
        bearer_token_file: /config/sa.token
        tls_config:
          ca_file: /config/ca.crt
         # insecure_skip_verify: true
    relabel_configs:
       - source_labels: [__address__] # 原标签
         regex: '(.*):10250'  # 匹配正则
         replacement: '${1}:9100'  #保留{1}，改9100端口
         target_label: __address__ 
         action: replace # 取代
```

### prometheus取得node中kubelet的cAdvisor
#### 访问kubelet提供的metrics服务
cAdvisor负责单节点内部的容器和节点资源使用统计，已经集成在 Kubelet 内部。
```bigquery
- job_name: 'k8s-kubelet'
    scheme: https
    bearer_token_file: /config/sa.token
    tls_config:
      ca_file: /config/ca.crt
    kubernetes_sd_configs:
      - api_server: https://<master ip>:6443/
        role: node # 基于节点，因为kubelet也是每个节点一个
        bearer_token_file: /config/sa.token
        tls_config:
          ca_file: /config/ca.crt
    relabel_configs: # 转换端口地址或ip地址
      - target_label: __address__
        replacement: <master ip>:6443
      - source_labels: [ __meta_kubernetes_node_name ]
        regex: '(.+)'
        replacement: '/api/v1/nodes/$1/proxy/metrics/cadvisor'
        target_label: __metrics_path__
        action: replace
```

### 安装Adapter
kubernetes主要通过两类 API 来获取资源使用指标：

resource metrics API:核心组件提供监控指标，如容器 CPU 、内存。 经典的实现就是metres-server。 HPA 也可以使用它来进行扩容。

custom metrics API：自定义指标。可以在HPA中使用自定义指标进行扩容。比较常用的就是prometheus-adapter

#### 安装
```bigquery
1. kubectl create custom-metrics
2. 修改custom-metrics-apiserver-deployment.yaml   修改prometheus对应的地址
3. # 进入adapter目录 
   kubectl apply -f .
4. 生成一个secret
    a. 进入k8s ca证书所在目录
    如果是kubeadm装的，默认目录在master主机的 /etc/kubernetes/pki

    b. 在/etc/kubernetes/pki目录执行以下几个命令:
    openssl genrsa -out serving.key 2048
    openssl req -new -key serving.key -out serving.csr -subj "/CN=serving"
    openssl x509 -req -in serving.csr -CA ./ca.crt -CAkey ./ca.key -CAcreateserial -out serving.crt -days 3650
    kubectl create secret generic cm-adapter-serving-certs --from-file=serving.crt=./serving.crt --from-file=serving.key -n   custom-metrics
5. 验证一下 kubectl get apiservice | grep custom-metrics  验证一下
[root@VM-0-16-centos ~]# kubectl get apiservice | grep custom-metrics
v1beta1.custom.metrics.k8s.io          custom-metrics/custom-metrics-apiserver   True        17h
v1beta1.external.metrics.k8s.io        custom-metrics/custom-metrics-apiserver   True        17h
v1beta2.custom.metrics.k8s.io          custom-metrics/custom-metrics-apiserver   True        17h
```

### prometheus监控自定义指标
1. 进入prometheus/custom_indicator目录
2. 编写事例 main.go文件
本次模拟用户访问量当做自定义指标
3. 在服务器上部署镜像 
```bigquery
docker run --rm -it -v /home/custom_indicator:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./usermetrics .
```
4. 编写deploy.yaml，启动服务
5. 修改prometheus.yml
法一：手动发现自定义指标
```bigquery
- job_name: 'xxx-metrics'
    static_configs:
    - targets: ['<deplyoment部署ip>:<service的端口>']
# 修改好都要记得重启
# curl -X POST http://<prometheus部署地址>:9090/-/reload    
```
法二：自动发现自定义指标
```bigquery
- job_name: 'xxxx-svc-auto'
    metrics_path: /metrics
    kubernetes_sd_configs:
    - api_server: https://<master ip>>:6443/
      role: service
      bearer_token_file: /config/sa.token
      tls_config:
        ca_file: /config/ca.crt
    relabel_configs: # 可以比对deploy.yaml中的Service
      - source_labels: [__meta_kubernetes_service_annotation_scrape]
        regex: true
        action: keep
      - source_labels: [__meta_kubernetes_service_annotation_nodeport]
        regex: '(.+)'
        replacement: '<deployment部署ip>:${1}'
        target_label: __address__
        action: replace
```
6. 需要修改adapter中的yaml文件，都在custom-indicator目录中
   custom-metrics-config-map.yaml
   custom-metrics-apiserver-deployment.yaml

### 自定义指标 HPA自动扩缩容
kubectl apply -f userhpa.yaml
```bigquery
[root@VM-0-16-centos ~]# kubectl get hpa
NAME      REFERENCE                TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
prodhpa   Deployment/usermetrics   0/8       1         3         1          13h
```
