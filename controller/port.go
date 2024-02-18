package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Port portt

type portt struct {
}

func (p *portt) PortTel(c *gin.Context) {
	portdata := new(service.PortData)
	if err := c.ShouldBindJSON(portdata); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	fmt.Println("接收到：", portdata)
	err := service.Port.TCPTelnet(portdata)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "TCP端口通信正常",
	})
}
