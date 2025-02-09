package main

import (
	"birth-info-service/config"
	"birth-info-service/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置 gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建请求处理器
	mailConfig := handlers.MailConfig{
		SMTPServer:   cfg.Mail.SMTPServer,
		SMTPPort:     cfg.Mail.SMTPPort,
		SMTPUsername: cfg.Mail.SMTPUsername,
		SMTPPassword: cfg.Mail.SMTPPassword,
	}

	handler := handlers.NewRequestHandler(
		cfg.API.Endpoint,
		mailConfig,
		cfg.API.DifyAPIKey,
	)

	// 初始化 Gin 引擎
	router := handlers.InitGin()

	// 配置路由
	router.POST("/info", handler.HandleBirthRequest)

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
