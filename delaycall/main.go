package delaycall

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Request 表示一个可能需要延迟处理的请求
type Request struct {
	UID        string // 用户唯一标识
	Payload    string // 请求载荷
	NeedsDelay bool   // 是否需要延迟
	id         int64
}

// delay 定义延迟时长为 100 毫秒
const delay = 100 * time.Millisecond

var (
	mu             sync.Mutex                      // 保护 activeDelayers 的互斥锁
	externalCh     = make(chan Request, 100)       // 外部请求通道
	activeDelayers = make(map[string]chan Request) // 存储正在延迟处理的用户通道
	checkCh        = make(chan int64, 1024)
)

// processor 从 externalCh 中读取请求并分发处理
func processor() {
	id := uuid.NewString()
	for req := range externalCh {
		mu.Lock()
		ch, delaying := activeDelayers[req.UID]
		if delaying {
			//fmt.Printf("[%s] DELAYER[%s]【%s】 延迟中...\n", time.Now().Format("15:04:05.000"), req.UID, req.Payload)
			// 用户正在延迟模式，将请求路由到对应的 delayer
			ch <- req
			//fmt.Printf("[%s] DELAYER[%s]【%s】 延迟发送\n", time.Now().Format("15:04:05.000"), req.UID, req.Payload)
			mu.Unlock()
			continue
		}
		if req.NeedsDelay {
			// 首次遇到需要延迟的请求，为该用户启动 delayer 协程
			ch = make(chan Request, 1000)
			activeDelayers[req.UID] = ch
			ch <- req
			mu.Unlock()
			go userDelayer(req.UID, ch)
			continue
		}
		// 普通请求：在主流程中顺序执行
		mu.Unlock()
		fmt.Printf("[%s] [%s] MAIN 处理 UID=%s Payload=%s\n", time.Now().Format("15:04:05.000"), id, req.UID, req.Payload)
		callService(req)
	}
}

// userDelayer 处理单个 UID 的延迟请求
func userDelayer(uid string, ch chan Request) {
	id := uuid.NewString()
	//fmt.Printf("[%s] [%s] DELAYER[%s] 启动\n", time.Now().Format("15:04:05.000"), id, uid)
	// 接收第一个请求
	req, ok := <-ch
	if !ok {
		mu.Lock()
		delete(activeDelayers, uid)
		close(ch)
		mu.Unlock()
		return
	}

	callService(req)
	// 处理请求
	fmt.Printf("[%s] [%s] DELAYER[%s] 处理 %s\n", time.Now().Format("15:04:05.000"), id, uid, req.Payload)

	// 等待新的请求或延迟超时
	timer := time.NewTimer(delay)
	for {
		select {
		case nextReq, ok := <-ch:
			// 收到新请求，且在延迟时长内
			timer.Stop()
			if !ok {
				return
			}

			callService(nextReq)
			// 再次延迟后处理下一请求
			fmt.Printf("[%s] [%s] DELAYER[%s] 处理 %s\n", time.Now().Format("15:04:05.000"), id, uid, nextReq.Payload)
			if nextReq.NeedsDelay {
				time.Sleep(delay)
			}

			timer.Reset(delay)
			continue

		case <-timer.C:
			mu.Lock()
			// 先检查一下ch中是否还有数据，如果有数据，本次不能关闭，重新延期
			if len(ch) > 0 {
				fmt.Printf("[%s] [%s] DELAYER[%s] 新增延期...\n", time.Now().Format("15:04:05.000"), id, uid)
				timer.Reset(delay)
				mu.Unlock()
				continue
			}

			delete(activeDelayers, uid)
			close(ch)
			mu.Unlock()
			// 超过延迟时长，退出延迟模式
			//fmt.Printf("[%s] [%s] DELAYER[%s] 延迟结束\n", time.Now().Format("15:04:05.000"), id, uid)
			timer.Stop()
			return
		}
	}
}

// callService 模拟外部服务调用
func callService(r Request) {
	//fmt.Printf("[%s] 调用服务 UID=%s Payload=%s\n",
	//	time.Now().Format("15:04:05.000"), r.UID, r.Payload)
	// 模拟执行耗时
	checkCh <- r.id
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func init() {
	nextId := int64(0)
	// 启动一个 goroutine 模拟外部服务调用
	go func() {
		for id := range checkCh {
			if nextId != id {
				fmt.Printf("[%s] [%d] 检测到 ID 不一致，请检查代码\n", time.Now().Format("15:04:05.000"), id)
				panic("")
			}

			nextId = id + 1
		}
	}()
}

// simulateRequests 随机生成 count 条 request，全部 UID=user1，
func simulateRequests(count int) {
	for i := 0; i < count; i++ {
		// 1/1000 概率需要延迟
		needsDelay := rand.Intn(100) == 0
		// 随机生成 payload，比如 task0 ~ task999
		payload := fmt.Sprintf("task%03d", i)

		externalCh <- Request{"user1", payload, needsDelay, int64(i)}

		// 随机 sleep 0–200ms
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	}
	close(externalCh)
}
