package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// 直连
type DirectDialer struct {
	Timeout time.Duration
}

func (d *DirectDialer) Dial(network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, d.Timeout)
}

// 代理
type Proxy struct {
	Addr string //ip:port
	// Type string //socks5/http默认socks5
}
type ProxyPool struct {
	Proxies []Proxy
	idx     int
	mu      sync.Mutex
}

func (p *ProxyPool) Next() Proxy {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Proxies) == 0 {
		return Proxy{}
	}
	proxy := p.Proxies[p.idx%len(p.Proxies)]
	p.idx++
	return proxy
}
func (p *ProxyPool) Dial(network, address string) (net.Conn, error) {
	maxRetries := 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		Proxy := p.Next()
		dialer, err := proxy.SOCKS5("tcp", Proxy.Addr, nil, &net.Dialer{
			Timeout: 2 * time.Second, // 代理连接超时（和直连保持一致）
		})
		if err != nil {
			lastErr = fmt.Errorf("代理%s创建失败: %w", Proxy.Addr, err)
			continue
		}
		conn, err := dialer.Dial(network, address)
		if err != nil {
			lastErr = fmt.Errorf("代理%s连接目标失败: %w", Proxy.Addr, err)
			continue
		}
		fmt.Println("✨ 使用代理：", Proxy.Addr)
		return conn, nil
		// 	dialer, err := proxy.SOCKS5("tcp", Proxy.Addr, nil, proxy.Direct)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	fmt.Println("✨ 使用代理：", Proxy.Addr)
		// 	return dialer.Dial(network, address)
		// }
	}
	return nil, fmt.Errorf("重试%d个代理均失败: %w", maxRetries, lastErr)
}
func LoadProxyPool(path string) ([]Proxy, error) {
	if path == "" {
		return nil, fmt.Errorf("proxy file path is empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open proxy file failed: %w", err)
	}
	defer f.Close()

	var proxies []Proxy
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		proxies = append(proxies, Proxy{Addr: line})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read proxy file failed: %w", err)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no valid proxies loaded from %s", path)
	}

	return proxies, nil
}
