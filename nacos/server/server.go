package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

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

	// 当前服务监听地址
	// 如果你的机器有多个网卡，可以使用 net.InterfaceAddrs() 获取真实 IP，本示例直接使用本机IP
	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("获取本机IP失败: %v", err)
	}
	port := 8080
	serviceName := "demo-service"

	// 注册服务到 Nacos
	instanceParam := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		Weight:      1.0,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	success, err := namingClient.RegisterInstance(instanceParam)
	if err != nil || !success {
		log.Fatalf("服务注册失败: %v", err)
	}
	log.Println("服务注册成功！")

	defer func() {
		_, err = namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			ServiceName: serviceName,
			Ip:          ip,
			Port:        uint64(port),
			GroupName:   "DEFAULT_GROUP",
			Ephemeral:   true,
		})
		if err != nil {
			log.Printf("注销服务实例失败: %v", err)
		} else {
			log.Println("服务实例注销成功")
		}
	}()

	// 启动 HTTP 服务器，提供 /hello 接口
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from server")
	})
	addr := fmt.Sprintf(":%d", port)
	log.Printf("HTTP 服务启动在 %s\n", addr)

	srv := &http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 等待信号（如 SIGINT 或 SIGTERM），然后开始优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("接收到关闭信号，开始注销服务并关闭HTTP服务...")

	// 注销服务实例
	_, err = namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: serviceName,
		Ip:          ip,
		Port:        uint64(port),
		GroupName:   "DEFAULT_GROUP",
		Ephemeral:   true,
	})
	if err != nil {
		log.Printf("注销服务实例失败: %v", err)
	} else {
		log.Println("服务实例注销成功")
	}

	// 优雅关闭 HTTP 服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP服务关闭失败: %v", err)
	}
	log.Println("HTTP服务关闭成功")
}

// getLocalIP 返回本机非环回的IP地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("未找到非环回IP地址")
}
