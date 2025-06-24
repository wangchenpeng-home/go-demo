package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Disable timestamp prefix in log output
	log.SetFlags(0)
	log.SetPrefix("")
	wsURL := "xxx"
	// 1. 建立 WebSocket 长连接
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	// 2. 启动一个 goroutine 专门读消息
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// msg 本身就是 []byte(JSON)，直接输出
			log.Println(string(msg))
			log.Println()
		}
	}()

	// 3. 发送登录
	loginPayload := map[string]interface{}{
		"op":   "login",
		"args": []interface{}{1000523071}, // 把这里替换成你自己的 accountId 或 token  1000523059
	}
	if err := conn.WriteJSON(loginPayload); err != nil {
		log.Fatalf("login write error: %v", err)
	}
	log.Printf("sent: %s", mustMarshal(loginPayload))

	// 4. 发送订阅
	subPnl := map[string]interface{}{
		"op":   "subscribe",
		"args": []interface{}{"gdfx.wallet", "gdfx.pnl", "gdfx.trade"},
	}

	if err := conn.WriteJSON(subPnl); err != nil {
		log.Fatalf("subscribe write error: %v", err)
	}
	log.Printf("sent: %s", mustMarshal(subPnl))

	// 5. 每 10 秒发送一次 ping
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		pingPayload := map[string]interface{}{
			"op":   "ping",
			"args": []interface{}{t.UnixNano() / int64(time.Millisecond)},
		}
		if err := conn.WriteJSON(pingPayload); err != nil {
			log.Printf("write ping error: %v", err)
			return
		}
		log.Printf("sent: %s", mustMarshal(pingPayload))
	}
}

// mustMarshal 是一个小助手，用于格式化打印 JSON payload
func mustMarshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
