# go_Scan
本项目是基于Go语言开发的高性能端口扫描器，支持并发扫描和代理池功能，帮助安全人员快速发现目标主机开放端口

## 功能特点

- 高性能并发端口扫描
- 支持tcp扫描代理池功能，提高扫描匿名性
- 轻量级设计，仅依赖Go官方库
- 简单易用的命令行接口

## 使用示例

基本扫描

```
#扫描单个主机的指定端口
go run main.go -u example.com -p 80,443,8080
#扫描单个主机的端口范围
go run main.go -u example.com -p 1-1024
#并发扫描
go run main.go -u example.com -p 1-1024 -w 100
#udp扫描
go run main.go -u example.com -p 1-1024 -s udp
```

使用代理池扫描：

```
go run main.go -u example.com -p 1-1024 -proxy -pf 1.txt
```

## 命令参数说明

| 参数   | 说明                                     | 示例                            |
| ------ | ---------------------------------------- | ------------------------------- |
| -u     | 目标主机/IP地址（仅支持单个）            | -u example.com                  |
| -p     | 端口（支持单个、多个或范围，用逗号分隔） | -p 80 或 -p 80,443 或 -p 1-1000 |
| -proxy | 开启代理池                               | -proxy                          |
| -py    | 代理池文件路径（相对于main.go的路径）    | -py proxy.txt                   |
| -w     | 并发数（默认不开启）                     | -w 100                          |
| -h     | 显示帮助信息                             | -h                              |

## 项目结构

```
websaomiao/           
├── internal/
 	└── cli/ 			# 命令行接口模块
│   ├── scanner/        # 核心扫描模块
│   └── util/           
	└── workpool/ 		#工作池模块
	└── main.go			#程序入口            

```

## 许可证

本项目采用Apache License 2.0许可证 - 详见[LICENSE ](https://github.com/zyhzheng500-maker/scan/blob/main/LICENSE)文件

## 作者

zyhsrc
