package routers

import (
	"GoScheduler"
	"GoScheduler/internal/modules/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"net/http"
)

func InitRouter(ctx *cli.Context) {
	router := gin.Default()

	// 注册通用中间件
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/", func(ctx *gin.Context) {
		content, err := GoScheduler.WebFS.ReadFile("index.html")
		if err != nil {

		}
		ctx.HTML(http.StatusOK, string(content), gin.H{})
	})

	host := parseHost(ctx)
	port := parsePort(ctx)
	router.Run(fmt.Sprintf("%s:%d", host, port))
}

// 注册专用中间件
func registerMiddleware(router *gin.Engine) {
	// IP 限制中间件

	// 登录认证中间件
}

// 解析端口
func parsePort(ctx *cli.Context) int {
	port := global.DefaultPort
	if ctx.IsSet("port") {
		port = ctx.Int("port")
	}
	if port <= 0 || port >= 65535 {
		port = global.DefaultPort
	}

	return port
}

func parseHost(ctx *cli.Context) string {
	if ctx.IsSet("host") {
		return ctx.String("host")
	}

	return "0.0.0.0"
}
