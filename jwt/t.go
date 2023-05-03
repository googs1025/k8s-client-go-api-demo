package main

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Enforcer *casbin.Enforcer

func init() {
	CasbinSetup()
}

// 初始化casbin
func CasbinSetup() {

	e := casbin.NewEnforcer("./rbac_models.conf")

	Enforcer = e

	// setup session store

	//return e
}

func Hello(c *gin.Context) {
	fmt.Println("Hello 接收到GET请求..")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
		"data": "Hello 接收到GET请求..",
	})
}

func main() {
	//获取router路由对象
	r := gin.New()

	// 增加policy
	r.GET("/api/v1/add", func(c *gin.Context) {
		fmt.Println("增加Policy")
		if ok := Enforcer.AddPolicy("admin", "*", "*"); !ok {
			fmt.Println("Policy已经存在")
		} else {
			fmt.Println("增加成功")
		}
	})
	//删除policy
	r.DELETE("/api/v1/delete", func(c *gin.Context) {
		fmt.Println("删除Policy")
		if ok := Enforcer.RemovePolicy("admin", "/api/v1/world", "GET"); !ok {
			fmt.Println("Policy不存在")
		} else {
			fmt.Println("删除成功")
		}
	})
	// 获取policy
	r.GET("/api/v1/get", func(c *gin.Context) {
		fmt.Println("查看policy")
		list := Enforcer.GetPolicy()
		for _, vlist := range list {
			for _, v := range vlist {
				fmt.Printf("value: %s, ", v)
			}
		}
	})

	// 使用自定义拦截器中间件
	r.Use(Authorize())

	// 创建请求
	r.GET("/api/v1/hello", Hello)
	r.Run(":9000") //参数为空 默认监听8080端口
}

// 拦截器

func Authorize() gin.HandlerFunc {

	return func(c *gin.Context) {
		//var e *casbin.Enforcer
		e := Enforcer

		//从DB加载策略
		//e.LoadPolicy()

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色 应该从db中读取
		sub := "admin"

		//判断策略中是否存在
		if ok := e.Enforce(sub, obj, act); ok {
			fmt.Println("恭喜您,权限验证通过")
			c.Next() // 进行下一步操作
		} else {
			fmt.Println("很遗憾,权限验证没有通过")
			c.Abort()
		}
	}
}
