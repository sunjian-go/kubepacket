package service

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"net"
	"time"
)

var Port portt

type portt struct {
}

type PortData struct {
	Ip      string `json:"ip"`
	TcpPort string `json:"tcpPort"`
}

func (p *portt) TCPTelnet(portdata *PortData) error {
	ip := portdata.Ip
	port := portdata.TcpPort

	// 创建 TCP 连接
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), 5*time.Second)
	if err != nil {
		logger.Error("TCP 连接失败:", err)
		return errors.New("TCP 连接失败:" + err.Error())
	}
	defer conn.Close()

	// 发送测试数据
	_, err = conn.Write([]byte("test data tcp"))
	if err != nil {
		logger.Error("发送TCP数据失败: ", err)
		return errors.New("发送TCP数据失败: " + err.Error())
	}

	fmt.Println(portdata.Ip + ":" + portdata.TcpPort + "：TCP端口通信正常")
	return nil
}
