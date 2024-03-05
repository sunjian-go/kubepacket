package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Pack pack

type pack struct {
}

var (
	process  *os.Process //当前命令的进程对象
	cmd      *exec.Cmd   //当前命令体
	status   = "0"       //是否上锁
	pcapName string      //保存的文件名
	timer    *time.Timer
)

type PackInfo struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	NetName string `json:"netName"`
	TimeOut string `json:"timeOut"`
}

// 获取所有网卡名
func (p *pack) GetAllInterface() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Error("获取接口列表失败" + err.Error())
		return nil, err
	}
	return interfaces, nil
}

// 生成随机名字并组合
func randName() string {
	namestr := "qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM1234567890"
	namearr := strings.Split(namestr, "")
	rand.Seed(time.Now().Unix()) //设置种子
	var newname []string
	for i := 0; i < 8; i++ {
		index := rand.Intn(len(namearr))
		newname = append(newname, namearr[index])
	}
	newName := strings.Join(newname, "")
	fmt.Println("新名字：", newName)
	return "pack-" + newName + ".pcap"
}

// 创建命令对象
func cmdobjFunc(packinfo *PackInfo) {
	if packinfo.Ip == "" && packinfo.Port == "" { //抓取所有包
		cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "-vv", "-nn", "-w", pcapName)
	} else if packinfo.Ip == "" && packinfo.Port != "" { //只指定端口抓包
		if strings.Contains(packinfo.Port, "-") { //只抓取端口范围
			cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "portrange", packinfo.Port, "-vv", "-nn", "-w", pcapName)
		} else {
			cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "port", packinfo.Port, "-vv", "-nn", "-w", pcapName)
		}
	} else if packinfo.Ip != "" && packinfo.Port != "" { //指定ip+端口抓包
		if strings.Contains(packinfo.Port, "-") { //指定ip+端口范围
			cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "portrange", packinfo.Port, "and", "host", packinfo.Ip, "-vv", "-nn", "-w", pcapName)
		} else {
			cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "port", packinfo.Port, "and", "host", packinfo.Ip, "-vv", "-nn", "-w", pcapName)
		}
	} else if packinfo.Ip != "" && packinfo.Port == "" { //只指定ip抓包
		cmd = exec.Command("tcpdump", "-i", packinfo.NetName, "host", packinfo.Ip, "-vv", "-nn", "-w", pcapName)
	}
}

func (p *pack) StartPacket(packinfo *PackInfo) error {
	//如果status == "1"说明已上锁，就直接返回
	if status == "1" {
		fmt.Println("当前已有抓包程序运行，请先停止当前抓包进程")
		return errors.New("当前已有抓包程序运行，请先停止当前抓包进程")
	}

	//如果status != "1"说明没上锁，就可以启动
	//获取随机name
	pcapName = randName()
	// 创建一个命令对象
	cmdobjFunc(packinfo)
	// 启动命令
	err := cmd.Start()
	if err != nil {
		fmt.Println("命令启动失败：", err.Error())
		//启动失败的话就解锁，等待下一次请求
		status = "0"
		return err
	}
	//启动成功的话就给当前抓包进程上锁
	status = "1"

	// 获取命令的进程对象
	process = cmd.Process

	fmt.Println("抓包程序已启动")
	//到这里代表已经启动抓包了

	go func() {
		//起携程去进行定时器
		timeout(packinfo.TimeOut)
	}()
	return nil
}

func (p *pack) StopPacket(cont *gin.Context) error {
	// 发送Ctrl+C信号给进程
	err := process.Signal(os.Interrupt)
	if err != nil {
		fmt.Println("发送停止抓包信号失败：", err.Error())
		return err
	}

	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Println("命令执行失败：", err.Error())
		return err
	}

	fmt.Println("抓包结束")

	//读取数据包写入响应体
	f, err := os.Open(pcapName)
	if err != nil {
		fmt.Println("打开文件失败：", err.Error())
		return err
	}

	// 获取文件的大小
	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("获取文件信息失败：%v\n", err)
		return err
	}

	cont.Header("Content-Type", "application/vnd.tcpdump.pcap")
	cont.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", pcapName)) //必须配置，指定下载名
	cont.Header("Content-Transfer-Encoding", "binary")

	//Content-Length: 响应中的主体内容的长度，以字节为单位。 必须配置，不然wireshark打开pcap文件会报错
	// "The capture file appears to be damaged or corrupt.(pcap: File has 3920856549-byte packet, bigger than maximum of262144)"
	cont.Header("Content-Length", fmt.Sprintf("%d", fi.Size()))
	fmt.Printf("文件大小：%d\n", fi.Size())

	n, err := io.Copy(cont.Writer, f)
	if err != nil {
		fmt.Println("写入响应体失败：", err.Error())
		return err
	}
	fmt.Println("写入字节：", n)
	//解锁抓包进程
	status = "0"
	//删除本地数据包
	os.Remove(pcapName)
	//结束定时器
	closeTimer()
	return nil
}

func timeout(timerr string) {
	//fmt.Println("超时时间：", timerr, "秒")
	sec, _ := strconv.Atoi(timerr)
	//创建定时器
	timer = time.NewTimer(time.Duration(sec) * time.Second)
	fmt.Println("定时器启动，时间：", timerr, "秒")
	//定时时间到
	<-timer.C
	if status == "1" {
		//说明还在抓包中
		//直接停掉抓包进程
		// 发送Ctrl+C信号给进程
		err := process.Signal(os.Interrupt)
		if err != nil {
			fmt.Println("发送停止抓包信号失败：", err.Error())
		}
		// 等待命令执行完成
		err = cmd.Wait()
		if err != nil {
			fmt.Println("命令执行失败：", err.Error())
		}
		fmt.Println("抓包结束")
		status = "0"
		//删除本地数据包
		os.Remove(pcapName)
		//结束定时器
		closeTimer()
	}
}

// 结束定时器
func closeTimer() {
	timer.Stop()
}
