// 端口逻辑
package util

type Port struct {
	PortChan chan int
}

func (p *Port) Close() {
	if p.PortChan != nil {
		close(p.PortChan) // 确保只关闭一次（避免重复关闭panic）
	}

}
