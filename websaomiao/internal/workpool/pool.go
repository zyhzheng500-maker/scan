// 工作池逻辑，并发扫描逻辑
package workpool

import (
	"fmt"
	"sync"
	"web/internal/scanner"
)

type WorkPool struct {
	//工人数量，自定义
	WorkerNum int
	wg        sync.WaitGroup
	Way       scanner.Scanner
	SharedRes scanner.Result
	ResMutex  sync.Mutex
}

func (wp *WorkPool) Start() {
	wp.wg.Add(wp.WorkerNum)
	//判断并发扫描方式
	switch v := wp.Way.(type) {
	case *scanner.Tcp:
		fmt.Println("并发启动,使用tcp扫描方式")

	default:
		fmt.Printf("并发启动失败,未知扫描方式%T", v)
		return
	}

	//启动工人池
	for i := 0; i < wp.WorkerNum; i++ {
		go func(i int) {
			defer wp.wg.Done()
			wp.Way.Scan(&wp.SharedRes, &wp.ResMutex)
		}(i)
	}
	wp.wg.Wait()
	fmt.Print("并发扫描完成\n")

}
