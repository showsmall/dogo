package api

import (
	"dogo/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Auth(c *gin.Context) {

	defer func() {
		token := c.GetHeader("X-Auth-Token")
		user, found := config.Cache.Get(token)
		if !found {
			Fail(c, 403, "no auth")
			c.Abort()
		} else {
			config.Cache.Set(token, user, time.Minute*time.Duration(30))
			c.Next()
		}
	}()
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", r)
			//debug.PrintStack()
			Fail(c, 0, fmt.Sprintf("%v", r))
			c.Abort()
		}
	}()
	c.Next()
}

func Cors(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "*")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	println(method, origin)

	//允许类型校验
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
		c.Abort()
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic info is: %v", err)
		}
	}()

	c.Next()
}
