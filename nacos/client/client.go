package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// 读取 Nacos 环境变量
	nacosHost := os.Getenv("NACOS_HOST")
	if nacosHost == "" {
		nacosHost = "127.0.0.1"
	}
	nacosPortStr := os.Getenv("NACOS_PORT")
	nacosPort := uint64(8848)
	if nacosPortStr != "" {
		if p, err := strconv.ParseUint(nacosPortStr, 10, 64); err == nil {
			nacosPort = p
		}
	}
	nacosUser := os.Getenv("NACOS_USER")
	nacosPassword := os.Getenv("NACOS_PASSWORD")

	// 配置 Nacos 服务器和客户端
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosHost,
			Port:   nacosPort,
		},
	}

	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		BeatInterval:        10000,
		NotLoadCacheAtStart: true,
		Username:            nacosUser,
		Password:            nacosPassword,
	}

	// 创建 Naming 客户端
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Fatalf("创建 Nacos Naming 客户端失败: %v", err)
	}

	// 根据服务名称获取服务实例列表
	serviceName := "demo-service"
	instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		HealthyOnly: true,
	})
	if err != nil {
		log.Fatalf("获取服务实例失败: %v", err)
	}

	if len(instances) == 0 {
		log.Fatalf("未找到服务实例")
	}

	// 选择第一个实例（可根据权重、健康状态进行选择）
	instance := instances[0]
	targetURL := fmt.Sprintf("http://%s:%d/hello", instance.Ip, instance.Port)
	log.Printf("调用服务 %s\n", targetURL)

	// 发起 HTTP 请求调用 /hello 接口
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Fatalf("调用服务失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}
	fmt.Printf("服务返回: %s\n", body)
}
