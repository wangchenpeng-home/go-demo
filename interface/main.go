package _interface

import "context"

type Module = *module

type module struct {
}

func (m Module) GetKey(ctx context.Context) string {
	return "key"
}

func (m Module) GetValue(ctx context.Context) string {
	return "value"
}

func (m Module) GetValueByKey(ctx context.Context, key string) string {
	getKey := m.GetKey(ctx)
	if getKey == "key" {
		return m.GetValue(ctx)
	}

	panic("not found key")
}
