package goline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	EndPoint  = "https://trialbot-api.line.me/v1/events"
	ToChannel = 1383378250
	EventType = "138311608800106203"
)

//BotRequestBody リクエスト内容を格納
type BotRequestBody struct {
	Request BotRequest `json:"Request"`
}

//BotRequest
type BotRequest struct {
	Result []BotResult `json:"result"`
}

type BotResult struct {
	From        string   `json:"from"`
	FromChannel string   `json:"fromChannel"`
	To          []string `json:"to"`
	ToChannel   string   `json:"toChannel"`
	EventType   string   `json:"eventType"`
	ID          string   `json:"id"`
	Content     Content  `json:"content"`
}

type Content struct {
	ID          string   `json:"id"`
	ContentType int      `json:"contentType"`
	From        string   `json:"from"`
	CreatedTime int      `json:"createdTime"`
	To          []string `json:"to"`
	ToType      int      `json:"toType"`
	Text        string   `json:"text"`
}
type SendRequest struct {
	To        []string `json:"to"`
	ToChannel int      `json:"toChannel"`
	EventType string   `json:"eventType"`
	Content   Content  `json:"content"`
}

type handler func([]BotResult) bool

type Server struct {
	router  *gin.Engine
	handler handler
}

func NewServer() *Server {
	s := &Server{}
	s.router = gin.New()
	s.router.Use(gin.Logger())
	return s
}

func (s *Server) SetHandler(h handler) {
	s.handler = h
}

func (s *Server) Run() {
	s.router.POST("/linebot/callback", func(c *gin.Context) {
		var botRequest BotRequest
		if c.Bind(&botRequest) == nil {
			c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("request convert error, request data is %+v", botRequest)})
			return
		}

		if ok := s.handler(botRequest.Result); !ok {
			c.JSON(http.StatusOK, gin.H{"status": "ng"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.router.Run(":" + os.Getenv("PORT"))
}

func Post(r SendRequest) (*http.Response, error) {
	b, _ := json.Marshal(r)
	req, _ := http.NewRequest(
		"POST",
		EndPoint,
		bytes.NewBuffer(b),
	)

	req = setHeader(req)

	proxyURL, _ := url.Parse(os.Getenv("FIXIE_URL"))
	client := &http.Client{
		Timeout:   time.Duration(15 * time.Second),
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}

	return client.Do(req)
}

func setHeader(req *http.Request) *http.Request {
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Line-ChannelID", os.Getenv("LINE_CHANNEL_ID"))
	req.Header.Add("X-Line-ChannelSecret", os.Getenv("LINE_CHANNEL_SECRET"))
	req.Header.Add("X-Line-Trusted-User-With-ACL", os.Getenv("LINE_CHANNEL_MID"))
	return req
}
