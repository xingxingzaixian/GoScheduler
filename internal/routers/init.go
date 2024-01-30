package routers

import (
	"GoScheduler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		content, err := GoScheduler.WebFS.ReadFile("index.html")
		if err != nil {

		}
		ctx.HTML(http.StatusOK, string(content), gin.H{})
	})
}
