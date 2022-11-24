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
### 部署方式
```bigquery
# 进入prometheus/state-metrics
kubectl apply -f .

# 部署好后可放问32280端口
```


