package routers

import (
	"GoScheduler/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

func InitRouter(host string, port int) {
	router := gin.Default()

	webFs, _ := fs.Sub(web.WebFS, "dist")
	router.StaticFS("/", http.FS(webFs))

	// 注册通用中间件
	router.Use(gin.Logger(), gin.Recovery())

	go func() {
		if err := router.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
			zap.S().Fatal(err)
		}
	}()
}

// 注册专用中间件
func registerMiddleware(router *gin.Engine) {
	// IP 限制中间件

	// 登录认证中间件
}
