package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
	"web/internal/util"
)

type Udp struct {
	Host string
	Port util.Port
}

func (u *Udp) Scan(sharedRes *Result, mu *sync.Mutex) error {
	var wg sync.WaitGroup
	openports := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for port := range openports {
			mu.Lock()
			sharedRes.Openports = append(sharedRes.Openports, port)
			mu.Unlock()
		}
	}()
	if u.Host == "" || u.Port.PortChan == nil {
		return fmt.Errorf("主机名Host或者端口Port不能为空")
	}
	for port := range u.Port.PortChan {
		fmt.Printf("正在扫描%v:%d的Udp端口\n", u.Host, port)
		conn, err := net.DialTimeout("udp", fmt.Sprintf("%v:%d", u.Host, port), 2*time.Second)

		if err != nil {
			continue
		}

		testData := []byte("Udp scan test")
		_, err = conn.Write(testData)
		if err != nil {
			continue
		}
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			continue
		} else {
			_ = n
			openports <- port
		}
		conn.Close()
	}
	close(openports)
	wg.Wait()
	return nil

}

// func (u *Udp) Scan(sharedRes *Result, mu *sync.Mutex) error {
// 	isOpen := true
// 	var wg sync.WaitGroup
// 	openports := make(chan int)
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for port := range openports {
// 			mu.Lock()
// 			sharedRes.Openports = append(sharedRes.Openports, port)
// 			mu.Unlock()
// 		}
// 	}()
// 	if u.Host == "" || u.Port.PortChan == nil {
// 		return fmt.Errorf("主机名Host或者端口Port不能为空")
// 	}
// 	retryNum := 1 // 重试1次（共2次尝试）
// 	timeout := 3 * time.Second

// 	for port := range u.Port.PortChan {
// 		fmt.Printf("正在扫描%v:%d的Udp端口\n", u.Host, port)
// 		for i := 0; i <= retryNum; i++ {
// 			targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", u.Host, port))
// 			if err != nil {
// 				isOpen = false
// 				break
// 			}
// 			conn, err := net.DialUDP("udp", nil, targetAddr)
// 			if err != nil {
// 				continue // 连接失败，进入下一次重试
// 			}
// 			conn.SetReadDeadline(time.Now().Add(timeout))
// 			conn.SetWriteDeadline(time.Now().Add(timeout))
// 			testData := []byte("Udp scan test")
// 			_, err = conn.Write(testData)
// 			if err != nil {
// 				conn.Close()
// 				continue
// 			}

// 			buf := make([]byte, 1024)
// 			n, readErr := conn.Read(buf)
// 			if readErr == nil {
// 				_ = n
// 				conn.Close()
// 				break
// 			} else {
// 				var opErr *net.OpError
// 				if errors.As(readErr, &opErr) {
// 					if opErr.Err == syscall.ECONNREFUSED || opErr.Err == syscall.ENETUNREACH {
// 						isOpen = false // 收到明确关闭信号，推翻默认假设
// 						conn.Close()
// 						break
// 					}
// 				}
// 			}
// 			conn.Close()
// 		}
// 		if isOpen {
// 			openports <- port
// 		}
// 	}

// 	close(openports)
// 	wg.Wait()
// 	return nil

// }
