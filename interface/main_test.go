package _interface

import (
	"context"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

// TestGetValueByKey_Normal 测试正常返回分支
func TestGetValueByKey_Normal(t *testing.T) {
	// 创建非 nil 的 module 实例
	m := Module(&module{})
	// 此时 m.GetKey 会返回 "key"，GetValueByKey 应返回 "value"
	result := m.GetValueByKey(context.Background(), "anyKey")
	if result != "value" {
		t.Fatalf("expected 'value', got '%s'", result)
	}
}

// TestGetValueByKey_Panic 测试触发 panic 的分支
func TestGetValueByKey_Panic(t *testing.T) {
	m := Module(&module{})
	// 使用 gomonkey 将 GetKey 方法补丁为返回非 "key" 的值，比如 "not_key"
	patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(m), "GetKey", []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{"not_key"},
		},
	})
	defer patches.Reset()

	// 捕获 panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic but did not occur")
		}
	}()

	// 调用时应触发 panic，因为 GetKey 返回 "not_key" 不满足 if 条件
	_ = m.GetValueByKey(context.Background(), "anyKey")
}
