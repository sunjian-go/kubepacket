package controller

import (
	"github.com/gin-gonic/gin"
	"main/service"
)

var Icmp icmp

type icmp struct {
}

func (i *icmp) PingFunc(c *gin.Context) {
	icmp := new(service.Icmpdata)
	if err := c.ShouldBindJSON(icmp); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}

	icmpresp, err := service.Icmp.PingFunc(icmp)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "测试ping成功",
		"data": icmpresp,
	})

}
