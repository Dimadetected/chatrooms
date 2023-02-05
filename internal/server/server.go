package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func New() *gin.Engine {
	wss := &webSocketSettings{
		mx:      &sync.Mutex{},
		clients: map[string]*webSocketClient{},
	}

	r := gin.New()
	r.GET("/echo", func(c *gin.Context) {
		initWsConn(wss, c)
	})

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	return r
}
