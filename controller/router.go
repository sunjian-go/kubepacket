package controller

import "github.com/gin-gonic/gin"

var Router router

type router struct {
}

func (router *router) Init(r *gin.Engine) {
	r.POST("/api/startPacket", Pack.StartPacket).
		POST("/api/stopPacket", Pack.StopPacket).
		GET("/api/interfaces", Pack.GetAllInterface).
		POST("/api/icmp", Icmp.PingFunc).
		POST("/api/port", Port.PortTel)
}
