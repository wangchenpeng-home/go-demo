// mock_ctx_mem_growth.go
package main

import (
	"context"
	"runtime"
)

// 定义一个空的 Span 和上下文 key，用于模拟真实场景
type Span struct {
	TraceID      string `json:"trace_id"`
	SpanID       string `json:"span_id"`
	ParentSpanID string `json:"parent_span_id"`
}

// 定义空类型用作context key
type empty struct{}

func main() {
	const N = 100000

	root := context.Background()
	var m1, m2 runtime.MemStats

	// 1. baseline：做一次 GC，读一次 stats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < N; i++ {
		_ = context.WithValue(root, empty{}, &empty{})
	}

	// 2. 批量 GC + stats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	delta := int64(m2.HeapAlloc) - int64(m1.HeapAlloc)
	println("WithValue alloc overhead total bytes:", delta)
	println("per call overhead (bytes):", delta/int64(N))
}
