package main

import (
	"github.com/gin-gonic/gin"
	"main/controller"
)

func main() {
	//创建路由引擎
	router := gin.Default()
	//初始化路由
	controller.Router.Init(router)

	router.Run("0.0.0.0:8888")
}
