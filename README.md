# 抓包插件
### 说明：
本程序是kubeutils项目的一个小插件，针对于k8s中node的功能研发，配合此插件，能使kubeutils项目实现远程抓包、icmp测试、端口测试（仅支持TCP）功能

# 构建可执行程序
```
go build -o kubepacket main.go
```

# 部署方式
使用二进制部署在你需要抓包、icmp测试以及端口测试的node上
```
chmod +x kubepacket && nohup ./kubepacket&
```
