package main

import (
	"fmt"
	"sync"
	"web/internal/cli"
	"web/internal/scanner"
	"web/internal/util"
	"web/internal/workpool"
)

func main() {
	var mu sync.Mutex
	var sharedRes scanner.Result
	config, err := cli.ParseCLI()
	if err != nil {
		fmt.Printf("命令行参数解析失败：%v\n", err)
		return
	}
	fmt.Printf("✅ 参数解析完成：\n")
	fmt.Printf("  目标主机：%s\n", config.Host)
	fmt.Printf("  待扫端口：%v\n", config.Ports)
	fmt.Printf("  并发数：%d\n", config.WorkerNum)
	fmt.Printf("  扫描类型：%s\n", config.ScanType)
	switch config.ScanType {
	case "tcp":

		portChan := make(chan int, len(config.Ports))
		for _, port := range config.Ports {
			portChan <- port
		}
		// close(portChan)
		portUtil := util.Port{PortChan: portChan}
		Tcpscanner := &scanner.Tcp{
			Host: config.Host,
			Port: portUtil,
		}
		if config.WorkerNum != 0 {
			fmt.Print("并发数不为0,开启并发扫描\n")
			pool := workpool.WorkPool{
				WorkerNum: config.WorkerNum,
				Way:       Tcpscanner,
				SharedRes: sharedRes,
			}
			go portUtil.Close()
			pool.Start()
			sharedRes = pool.SharedRes //值传递，最后需要将结果赋值回去
		} else {
			fmt.Println("并发数为0,不开启并发扫描")
			go portUtil.Close()
			Tcpscanner.Scan(&sharedRes, &mu)

		}
	case "udp":
		portChan := make(chan int, len(config.Ports))
		for _, port := range config.Ports {
			portChan <- port
		}
		portUtil := util.Port{PortChan: portChan}
		Udpsacnner := &scanner.Udp{
			Host: config.Host,
			Port: portUtil,
		}
		if config.WorkerNum != 0 {
			fmt.Print("并发数不为0,开启并发扫描\n")
			pool := workpool.WorkPool{
				WorkerNum: config.WorkerNum,
				Way:       Udpsacnner,
				SharedRes: sharedRes,
			}
			go portUtil.Close()
			pool.Start()
			sharedRes = pool.SharedRes //值传递，最后需要将结果赋值回去
		} else {
			fmt.Println("并发数为0,不开启并发扫描")
			go portUtil.Close()
			Udpsacnner.Scan(&sharedRes, &mu)
		}
	}
	if len(sharedRes.Openports) == 0 {
		fmt.Print("扫描完成，未发现开放端口\n")
		return
	}
	fmt.Printf("扫描完成，开放端口：%v\n", sharedRes.Openports)

}
