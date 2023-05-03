package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
	使用prometheus自定义接口指标
*/

func init() {
	prometheus.MustRegister(userVisit)
}

var userVisit = prometheus.NewCounterVec(
	// 指标采用计数器模式：
	prometheus.CounterOpts{
		Name: "jiang_user_visit", // 指标名称(prometheus)查询的时候会用到。
	},
	[]string{"userid"},
)

func main() {

	r := gin.New()
	// 自定义的业务接口：模拟用户的访问量
	r.GET("/user/visit", func(c *gin.Context) {

		userStr := c.Query("userid")

		fmt.Printf("the user is %s\n", userStr)

		userVisit.With(prometheus.Labels{
			"userid": userStr,
		}).Inc() // 访问量指标就增加一次

		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	// 要填写，prometheus才能抓取到
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8089")

}
