package delaycall

import (
	"math/rand"
	"testing"
	"time"
)

func Test_userDelayer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	// 启动主处理流程
	go processor()

	// 模拟请求输入
	simulateRequests(1000000)

	// 让示例运行一段时间
	time.Sleep(10000 * time.Second)
}
