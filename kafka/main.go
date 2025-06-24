package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

// 實現 SCRAMClient 接口
type SCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *SCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *SCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *SCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func main() {
	// 設定 Sarama 日誌
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	// 創建配置
	config := sarama.NewConfig()

	// 設置 SASL/SCRAM 認證
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512 // 或 SASLTypeSCRAMSHA256
	config.Net.SASL.User = "bybit"                         // MSK 用戶名
	config.Net.SASL.Password = "LI+b09|Wi[29lIiy=2}+"
	config.Producer.Return.Successes = true
	// 設置 SCRAM 客戶端生成函數
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &SCRAMClient{
			HashGeneratorFcn: scram.SHA512, // 或 scram.SHA256
		}
	}

	// 啟用 TLS
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		MinVersion: tls.VersionTLS12,
		// 生產環境應該使用正確的證書驗證
		// InsecureSkipVerify: true, // 僅開發測試環境使用
	}

	// 設置 Kafka 版本，根據你的 MSK 版本調整
	config.Version = sarama.V2_8_1_0

	// MSK Broker 地址
	brokers := []string{
		"b-3-public.bybitdemomskkafka.75z2xt.c4.kafka.ap-northeast-1.amazonaws.com:9196",
		"b-1-public.bybitdemomskkafka.75z2xt.c4.kafka.ap-northeast-1.amazonaws.com:9196",
		"b-2-public.bybitdemomskkafka.75z2xt.c4.kafka.ap-northeast-1.amazonaws.com:9196",
	}

	// 嘗試建立連接
	client, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	fmt.Println("成功連接到 MSK Kafka 叢集!")

	// 接下來可以創建生產者或消費者
	// producer, err := sarama.NewSyncProducerFromClient(client)
	// consumer, err := sarama.NewConsumerFromClient(client)
	// // Topic 設定
	//topicName := "demo-topic"
	//numPartitions := 3
	//replicationFactor := 3

	// Topic 詳細配置 (可選)
	//topicDetail := &sarama.TopicDetail{
	//	NumPartitions:     int32(numPartitions),
	//	ReplicationFactor: int16(replicationFactor),
	//	ConfigEntries:     map[string]*string{
	//		// 可選配置項，例如:
	//		// "cleanup.policy": stringPtr("compact"),
	//		// "retention.ms": stringPtr("604800000"), // 7 天
	//	},
	//}

	// 建立 Topic
	//err = client.CreateTopic(topicName, topicDetail, false)
	//if err != nil {
	//	// 檢查是否是因為 topic 已存在而失敗
	//	if err == sarama.ErrTopicAlreadyExists {
	//		fmt.Printf("Topic '%s' 已經存在\n", topicName)
	//	} else {
	//		log.Fatalf("建立 Topic 失敗: %v", err)
	//	}
	//} else {
	//	fmt.Printf("Topic '%s' 成功建立，分區數: %d，複製因子: %d\n",
	//		topicName, numPartitions, replicationFactor)
	//}

	// 创建同步生产者（用于测试连接）
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("关闭生产者失败: %v", err)
		}
	}()

	// 测试发送消息
	topic := "test-topic1" // 替换为实际存在的Topic
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello Kafka from Go!"),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	fmt.Printf("消息发送成功！Partition: %d, Offset: %d\n", partition, offset)

	//client.DeleteTopic(topicName)

	// // 列出所有 topics 以驗證 (可選)
	// topics, err := cli.ListTopics()
	// if err != nil {
	// 	log.Fatalf("列出 Topics 失敗: %v", err)
	// }

	// fmt.Println("現有的 Topics:")
	// for name, detail := range topics {
	// 	fmt.Printf("- %s (分區數: %d, 複製因子: %d)\n",
	// 		name, detail.NumPartitions, detail.ReplicationFactor)
	// }
}
