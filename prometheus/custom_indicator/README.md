### prometheus自定义指标实践




Counter(计数器): 表示 一个单调递增的指标数据，请求次数、错误数量等等

Gauge(计量器):代表一个可以任意变化的指标数据， 可增可减 .场景有：协程数量、CPU、Memory 、业务队列的数量 等等

Histogram(累积直方图)
主要是样本观测数据,在一段时间范围内对数据进行采样. 譬如请求持续时间或响应大小等等 。这点往往可以配合链路追踪系统（譬如之前讲到过jaeger来使用）

Summary (摘要统计)
和 直方图类似，也是样本观测。但是它提供了样本值的分位数、所有样本值的大小总和、样本总量

### prometheus.yml文件中的配置方法：
```bigquery
[ source_labels: '[' <labelname> [, ...] ']' ] # 源标签从现有标签中选择值。
[ regex: <regex> | default = (.*) ] # 与提取的值匹配的正则表达式。
[ target_label: <labelname> ] # 被替换的标签。
[ replacement: <string> | default = $1 ] # 替换值
[ action: <relabel_action> | default = replace ]匹配执行的操作。 （ replace 、keep、drop、labelmap、labeldrop）
```

### keep和drop的作用
当action设置为keep时，Prometheus会丢弃source_labels的值中没有匹配到regex正则表达式内容的Target实例，

而当action设置为drop时，则会丢弃那些source_labels的值匹配到regex正则表达式内容的Target实例

