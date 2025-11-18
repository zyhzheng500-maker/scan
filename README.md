# scan
本项目是基于go语言实现的端口扫描器，支持并发扫描，帮助安全人员快速发现目标主机开放端口



# 使用示例

基本扫描

```
#扫描单个主机的指定端口
go run main.go --host example.com --ports 80,443,8080
#扫描单个主机的端口范围
go run main.go --host example.com --ports 1-1024
#并发扫描
go run main.go --host example.com --ports 1-1024 --worker 100
```

