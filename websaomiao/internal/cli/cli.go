// 命令行指令
package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLIConfig struct {
	Host         string
	Ports        []int
	WorkerNum    int
	ScanType     string
	UseProxyPool bool
	ProxyFile    string
}

func ParseCLI() (CLIConfig, error) {
	var (
		host     = flag.String("u", "", "目标主机/IP(必选,示例:127.0.0.1 或 baidu.com)")
		portsStr = flag.String("p", "", "待扫描端口（必选,单个端口如80,多个端口用逗号分隔如80,443,也可以用1-1024来表示端口范围,也用逗号隔开)") //暂时先只用逗号，后面用1-1000这种范围的再扩展
		worker   = flag.Int("w", 0, "并发数(可选,默认3)")
		scanType = flag.String("s", "tcp", "扫描类型(默认tcp,当前支持tcp,udp)")
		// 是否启用代理池
		useProxyPool = flag.Bool("proxy", false, "是否启用代理池")
		// 代理列表文件
		proxyFile = flag.String("pf", "", "代理池文件，每行一个 host:port")
	)

	// 2. 解析用户输入的命令行参数（必须调用flag.Parse()，否则参数无法生效）
	flag.Parse()
	if flag.NFlag() == 0 || (flag.Parsed() && (*host == "" || *portsStr == "")) {
		// 打印程序使用说明（告诉用户如何输入参数）
		fmt.Println("使用说明：./scan-tool [参数]")
		// 打印所有参数的默认值和帮助说明（自动生成，基于定义时的第三个参数）
		flag.PrintDefaults()
		// 若必选参数缺失，返回错误信息（引导用户补充参数）
		if *host == "" || *portsStr == "" {
			return CLIConfig{}, fmt.Errorf("错误：--host 和 --ports 为必选参数，请重新输入")
		}
		// 正常显示帮助后，主动退出程序（避免继续执行）
		os.Exit(0)
	}

	// 4. 校验并解析端口参数（用户输入字符串→转int切片）
	ports, err := parsePorts(*portsStr)
	if err != nil {
		return CLIConfig{}, fmt.Errorf("端口解析失败：%v", err)
	}

	// 5. 校验其他参数合法性
	if err := validateConfig(*host, ports, *worker, *scanType, *useProxyPool, *proxyFile); err != nil {
		return CLIConfig{}, fmt.Errorf("参数校验失败：%v", err)
	}

	// 6. 返回结构化参数（将flag的指针值转为具体值）
	return CLIConfig{
		Host:         *host,
		Ports:        ports,
		WorkerNum:    *worker,
		ScanType:     *scanType,
		UseProxyPool: *useProxyPool,
		ProxyFile:    *proxyFile,
	}, nil
}

// cli/cli.go 补充端口解析函数
func parsePorts(portsStr string) ([]int, error) {
	var ports []int
	portSet := make(map[int]struct{}) //用于去重

	portStrList := strings.Split(portsStr, ",")

	for _, pStr := range portStrList {
		pStr = strings.TrimSpace(pStr)
		if pStr == "" {
			continue
		}
		if strings.Contains(pStr, "-") {
			rangeParts := strings.Split(pStr, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("端口范围格式错误:%s", pStr)
			}
			startStr := strings.TrimSpace(rangeParts[0])
			endStr := strings.TrimSpace(rangeParts[1])
			startPort, err1 := strconv.Atoi(startStr)
			endPort, err2 := strconv.Atoi(endStr)
			if err1 != nil {
				return nil, fmt.Errorf("端口范围起始值无效：%s(必须是数字)", startStr)
			}
			if err2 != nil {
				return nil, fmt.Errorf("端口范围结束值无效：%s(必须是数字)", endStr)
			}
			if startPort > endPort {
				return nil, fmt.Errorf("端口范围无效：%s(起始端口%d > 结束端口%d)", pStr, startPort, endPort)
			}
			if startPort < 1 || endPort > 65535 {
				return nil, fmt.Errorf("端口范围超出合法范围1-65535:%s", pStr)
			}
			for port := startPort; port <= endPort; port++ {
				if _, exists := portSet[port]; !exists {
					portSet[port] = struct{}{} //空结构体实例
					ports = append(ports, port)
				}
			}
		} else {
			port, err := strconv.Atoi(pStr)
			if err != nil {
				return nil, fmt.Errorf("端口格式无效:%s(必须是数字)", pStr)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("端口超出合法范围1-65535:%d", port)
			}
			if _, exists := portSet[port]; !exists {
				portSet[port] = struct{}{}
				ports = append(ports, port)
			}
		}

	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("未解析到任何有效端口,请检查--ports参数")
	}
	return ports, nil
}

// cli/cli.go 补充参数校验函数
func validateConfig(host string, ports []int, worker int, scanType string, useProxyPool bool, proxyFile string) error {
	// 校验主机：非空（已在ParseCLI中初步判断，这里二次确认）
	if host == "" {
		return fmt.Errorf("目标主机--host不能为空")
	}

	// 校验端口：每个端口必须在1-65535范围
	for _, port := range ports {
		if port < 1 || port > 65535 {
			return fmt.Errorf("端口%d不合法(合法范围1-65535)", port)
		}
	}

	// 校验并发数：≥1（避免工作池无法启动）
	if worker < 0 {
		return fmt.Errorf("并发数--worker必须≥1(当前输入:%d)", worker)
	}

	// 校验扫描类型：当前仅支持tcp（后续扩展UDP时只需加case）
	if scanType != "tcp" && scanType != "udp" {
		return fmt.Errorf("不支持的扫描类型--scan-type:%s(当前仅支持tcp,udp)", scanType)
	}
	if useProxyPool && proxyFile == "" {
		return fmt.Errorf("启用代理池(--proxy)时，必须通过--proxy-file指定代理文件路径")
	}
	return nil
}
