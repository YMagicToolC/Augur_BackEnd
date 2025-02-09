package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	API       APIConfig       `yaml:"api"`
	Mail      MailConfig      `yaml:"mail"`
	Logging   LoggingConfig   `yaml:"logging"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Timeouts  TimeoutConfig   `yaml:"timeouts"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	Mode string `yaml:"mode"`
}

type APIConfig struct {
	Endpoint   string `yaml:"endpoint"`
	DifyAPIKey string `yaml:"dify_api_key"`
}

type MailConfig struct {
	SMTPServer   string `yaml:"smtp_server"`
	SMTPPort     int    `yaml:"smtp_port"`
	SMTPUsername string `yaml:"smtp_username"`
	SMTPPassword string `yaml:"smtp_password"`
	SenderName   string `yaml:"sender_name"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type RateLimitConfig struct {
	Enabled           bool    `yaml:"enabled"`
	RequestsPerSecond float64 `yaml:"requests_per_second"`
}

type TimeoutConfig struct {
	APIRequest int `yaml:"api_request"`
	EmailSend  int `yaml:"email_send"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// 读取配置文件
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件错误: %v", err)
	}

	// 解析YAML
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件错误: %v", err)
	}

	// 验证必要的配置
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	// API配置验证
	if c.API.Endpoint == "" {
		return fmt.Errorf("API endpoint 不能为空")
	}
	if c.API.DifyAPIKey == "" {
		return fmt.Errorf("Dify API key 不能为空")
	}

	// 邮件配置验证
	if c.Mail.SMTPServer == "" || c.Mail.SMTPUsername == "" || c.Mail.SMTPPassword == "" {
		return fmt.Errorf("邮件配置不完整")
	}

	return nil
}
