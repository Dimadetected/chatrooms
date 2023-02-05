package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type webSocketSettings struct {
	clients map[string]*webSocketClient
	mx      *sync.Mutex
}

type webSocketClient struct {
	name string
	conn *websocket.Conn
}

type message struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func initWsConn(wss *webSocketSettings, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln(err)
	}

	wss.AddClient(conn, c.Request.URL.Query().Get("name"))
}

func (wss *webSocketSettings) AddClient(conn *websocket.Conn, name string) {
	if name == "" {
		name = "Новый пользователь"
	}

	newClient := &webSocketClient{
		name: name,
		conn: conn,
	}

	wss.mx.Lock()

	stringUUID := uuid.New().String()
	wss.clients[stringUUID] = newClient
	wss.mx.Unlock()

	wss.CheckClientMessages(newClient, stringUUID)
}
func (wss *webSocketSettings) CheckClientMessages(client *webSocketClient, stringUUID string) {

	for {
		msgType, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseProtocolError) {
				wss.mx.Lock()
				delete(wss.clients, stringUUID)
				wss.mx.Unlock()

				return
			}
			log.Fatalln("read message", err)
		}

		wss.SendMessagesToClients(msgType, msg)
	}
}

func (wss *webSocketSettings) SendMessagesToClients(msgType int, msg []byte) {
	for _, client := range wss.clients {
		jsonMsg, err := json.Marshal(message{
			Data: string(msg),
			Name: client.name,
		})
		if err != nil {
			log.Fatalln("marshall msg", err)
		}

		if err := client.conn.WriteMessage(msgType, jsonMsg); err != nil {
			log.Fatalln("write msg to socket", err)
		}
	}
}
