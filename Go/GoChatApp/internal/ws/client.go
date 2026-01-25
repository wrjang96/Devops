package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second // 클라이언트가 서버 ping에 대해 pong 응답을 60초 안에 줘야 정상
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 고쳐야 할 부분
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan Message // Message Queue
	room     string
	username string
}

func (c *Client) readPump() { // client가 보내는 메시지를 계속 읽음
	defer func() {
		c.hub.unregister <- Subscription{Client: c, Room: c.room}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := Message{
			Room:    c.room,
			Sender:  c.username,
			Content: string(messageBytes),
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, c *gin.Context) {
	room := c.Query("room")
	username := c.Query("username")

	if room == "" {
		room = "general"
	}
	if username == "" {
		username = "Anonymous"
	}

	// Check limit
	if hub.IsRoomFull(room) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Room is full (max 50 users)"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, 256),
		room:     room,
		username: username,
	}
	client.hub.register <- Subscription{Client: client, Room: room}

	go client.writePump()
	go client.readPump()
}
