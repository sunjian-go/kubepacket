package service

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/wonderivan/logger"
	"os/user"
	"strconv"
	"time"
)

var Icmp icmp

type icmp struct {
}

type Icmpdata struct {
	Ip      string `json:"ip"`
	TimeOut string `json:"timeOut"` //超时秒
	Count   string `json:"count"`   //数据包数量
}

type IcmpResp struct {
	Sent int     `json:"sent"`
	Recv int     `json:"recv"`
	Loss float64 `json:"loss"`
	Min  string  `json:"min"`
	Avg  string  `json:"avg"`
	Max  string  `json:"max"`
}

// ping方法
func (i *icmp) PingFunc(icmpdata *Icmpdata) (*IcmpResp, error) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("无法获取当前用户信息：", err)
		return nil, err
	}

	fmt.Println("当前用户ID：", currentUser.Uid)

	//超时时间（与下面ICMP请求数量对应）
	timenum, _ := strconv.Atoi(icmpdata.TimeOut)
	timeout := time.Second * time.Duration(timenum)

	pinger, err := ping.NewPinger(icmpdata.Ip)
	if err != nil {
		logger.Error("创建ping对象失败")
		return nil, err
	}
	pinger.SetPrivileged(true)               //特权运行。必须设置：否则运行会出现：socket: permission denied
	pinger.Timeout = timeout                 //设置超时时间
	count, _ := strconv.Atoi(icmpdata.Count) //获取icmp请求数量
	pinger.Count = count                     //发送3个ICMP请求，一秒发一个

	err = pinger.Run()
	if err != nil {
		logger.Error("Ping失败:", err)
		return nil, err
	}

	//通过Statistics方法获取Ping的统计信息
	stats := pinger.Statistics()
	fmt.Printf("Ping 状态信息: 发送 = %d, 接收 = %d, 丢包 = %v\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("往返时间的最小值/平均值/最大值): %v/%v/%v\n", stats.MinRtt.String(), stats.AvgRtt.String(), stats.MaxRtt.String())

	//组装数据返回
	icmpresp := &IcmpResp{
		Sent: stats.PacketsSent,
		Recv: stats.PacketsRecv,
		Loss: stats.PacketLoss,
		Min:  stats.MinRtt.String(),
		Avg:  stats.AvgRtt.String(),
		Max:  stats.MaxRtt.String(),
	}
	return icmpresp, nil
}
