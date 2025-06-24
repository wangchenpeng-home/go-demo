package json

import (
	"encoding/json"
	"testing"

	sonic "github.com/bytedance/sonic"
	jsoniter "github.com/json-iterator/go"
)

// Order 定义了订单的结构体，与给定 JSON 数据对应
type Order struct {
	Account        int64  `json:"account"`
	PositionId     int64  `json:"positionId"`
	OrderType      int    `json:"orderType"`
	Symbol         string `json:"symbol"`
	Volume         string `json:"volume"`
	OpenTime       int64  `json:"openTime"`
	OpenPrice      string `json:"openPrice"`
	StopLimitPrice string `json:"stopLimitPrice"`
	TP             string `json:"tp"`
	SL             string `json:"sl"`
	Commission     string `json:"commission"`
	Taxes          string `json:"taxes"`
	Swap           string `json:"swap"`
	Profit         string `json:"profit"`
	TotalProfit    string `json:"totalProfit"`
	Comment        string `json:"comment"`
}

func TestB(t *testing.T) {

}

// 用于测试的 JSON 数据，实际使用时可根据需要构造更复杂或更大的样本数据
var jsonData = []byte(`{
	"account": 922801,
	"positionId": 296373361,
	"orderType": 0,
	"symbol": "BTCUSD",
	"volume": "0.01",
	"openTime": 1743674672,
	"openPrice": "87319.54",
	"stopLimitPrice": "0.00",
	"tp": "0.00",
	"sl": "0.00",
	"commission": "0.00",
	"taxes": "0.00",
	"swap": "0.00",
	"profit": "-0.03",
	"totalProfit": "-0.03",
	"comment": ""
}`)

// BenchmarkEncodingJSON 使用标准库 encoding/json 解析 JSON 的基准测试
func BenchmarkEncodingJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		order := &Order{}
		if err := json.Unmarshal(jsonData, &order); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter 使用第三方 jsoniter 解析 JSON 的基准测试
func BenchmarkJsoniter(b *testing.B) {
	// 采用与标准库兼容的配置
	var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		order := &Order{}
		if err := jsoni.Unmarshal(jsonData, &order); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSonic uses the sonic library for JSON parsing benchmark
func BenchmarkSonic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		order := &Order{}
		if err := sonic.Unmarshal(jsonData, &order); err != nil {
			b.Fatal(err)
		}
	}
}
