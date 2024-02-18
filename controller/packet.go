package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/service"
)

var Pack packet

type packet struct {
}

// 启动抓包进程
func (p *packet) StartPacket(c *gin.Context) {
	packinfo := new(service.PackInfo)
	if err := c.ShouldBindJSON(packinfo); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	fmt.Println("需要抓包的数据为：", packinfo)
	err := service.Pack.StartPacket(packinfo)
	if err != nil {
		fmt.Println("packet: ", err.Error())
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "启动抓包程序成功",
	})
}

// 停止抓包
func (p *packet) StopPacket(c *gin.Context) {
	err := service.Pack.StopPacket(c)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
}

// 获取所有网卡信息
func (p *packet) GetAllInterface(c *gin.Context) {
	interfaces, err := service.Pack.GetAllInterface()
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取接口列表成功",
		"data": interfaces,
	})
}
