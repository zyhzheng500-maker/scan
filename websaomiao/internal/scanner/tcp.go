// tcp扫描方式
package scanner

import (
	"fmt"
	"sync"
	"web/internal/util"
)

type Tcp struct {
	Host   string
	Port   util.Port
	Dialer Dialer
}

// tcp扫描功能
func (t *Tcp) Scan(sharedRes *Result, mu *sync.Mutex) error {

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
	if t.Host == "" || t.Port.PortChan == nil {
		return fmt.Errorf("主机名Host或者端口Port不能为空")
	}
	//不用管端口是单个还是多个，统一通过通道进行处理，只需要分开是否并发扫描即可
	for port := range t.Port.PortChan {
		fmt.Printf("正在扫描%v:%d的TCP端口\n", t.Host, port)

		conn, err := t.Dialer.Dial("tcp", fmt.Sprintf("%v:%d", t.Host, port))
		if err != nil {
			continue
		}
		openports <- port
		conn.Close()
	}
	close(openports)
	wg.Wait()

	return nil

}
