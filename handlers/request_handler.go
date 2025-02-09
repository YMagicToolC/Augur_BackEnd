package handlers

import (
	"birth-info-service/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

var log = logrus.New()

func init() {
	// 设置日志格式为JSON
	log.SetFormatter(&logrus.JSONFormatter{})
	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)

	// 打开或创建日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("无法将日志记录到文件，使用默认的标准错误输出")
	}
}

type RequestHandler struct {
	apiEndpoint string
	mailConfig  MailConfig
	difyAPIKey  string
	httpClient  *http.Client
}

type MailConfig struct {
	SMTPServer   string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

func NewRequestHandler(apiEndpoint string, mailConfig MailConfig, difyAPIKey string) *RequestHandler {
	if apiEndpoint == "" || difyAPIKey == "" {
		panic("apiEndpoint 和 difyAPIKey 不能为空")
	}
	if mailConfig.SMTPServer == "" || mailConfig.SMTPUsername == "" {
		panic("邮件配置不完整")
	}

	log.Info("RequestHandler initialized")

	return &RequestHandler{
		apiEndpoint: apiEndpoint,
		mailConfig:  mailConfig,
		difyAPIKey:  difyAPIKey,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (h *RequestHandler) HandleBirthRequest(c *gin.Context) {
	var request models.APIRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.WithFields(logrus.Fields{
		"request": request,
	}).Info("Received birth request")

	// 创建API请求

	// 异步发送到API并处理响应
	go h.sendToAPIAndNotify(request)

	c.JSON(http.StatusOK, gin.H{
		"contact": "",
		"message": "请求已接收，正在处理",
	})
}

func (h *RequestHandler) sendToAPIAndNotify(request models.APIRequest) {

	// 构建新的请求体
	requestBody := map[string]interface{}{
		"inputs":        request, // 假设 request 已经是符合 API 要求的格式
		"response_mode": "blocking",
		"user":          "yacoservice",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("序列化请求失败: %v", err)
		return
	}

	req, err := http.NewRequest("POST", h.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.difyAPIKey))

	log.WithFields(logrus.Fields{
		"apiRequest": requestBody,
	}).Info("Sending request to API")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API 响应错误，状态码: %d", resp.StatusCode)
		return
	}

	var apiResponse models.APIResponseData
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("解析响应失败: %v", err)
		return
	}

	if err := h.sendEmail(request.Contact, apiResponse.Data.Outputs); err != nil {
		log.Printf("发送邮件失败: %v", err)
		return
	}
}

func (h *RequestHandler) sendEmail(email string, response models.APIOutput) error {
	m := mail.NewMessage()
	m.SetHeader("From", h.mailConfig.SMTPUsername)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "出生信息查询结果")

	// 构建邮件内容
	body := fmt.Sprintf(`
        <p>尊敬的用户：</p>
        <p>您的结果如下：</p>
        <p>详细结果：%s </p>
		  
        <p>如有任何问题，请随时与我们联系。</p>
        <p>此致</p>
    `, response.Message)

	m.SetBody("text/html", body)

	d := mail.NewDialer(
		h.mailConfig.SMTPServer,
		h.mailConfig.SMTPPort,
		h.mailConfig.SMTPUsername,
		h.mailConfig.SMTPPassword,
	)

	log.WithFields(logrus.Fields{
		"email":    email,
		"response": response,
	}).Info("Sending email")

	if err := d.DialAndSend(m); err != nil {
		log.WithError(err).Error("Failed to send email")
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	return nil
}

func InitGin() *gin.Engine {
	r := gin.Default()

	// 使用 gin-contrib/cors 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有域名
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 方法 1：设置受信任的代理
	// r.SetTrustedProxies([]string{"192.168.1.2", "192.168.1.3"})

	// 或者 方法 2：如果在本地开发环境，可以禁用代理信任
	r.SetTrustedProxies(nil)

	return r
}
