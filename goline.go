package goline

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
