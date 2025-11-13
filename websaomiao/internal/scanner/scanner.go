package scanner

import (
	"sync"
)

// 扫描接口，以便后续添加更多扫描方式
type Scanner interface {
	Scan(sharedRes *Result, mu *sync.Mutex) error
}
type Result struct {
	Openports []int
}
