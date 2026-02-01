package ws

import (
	"log"
	"time"

	"go-chat-mem/internal/store"

	"github.com/google/uuid"
)

type Message struct {
	Room    string `json:"room"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Subscription struct {
	Client *Client
	Room   string
}

type RoomStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type LimitCheck struct {
	Room   string
	Result chan bool
}

type Hub struct {
	rooms map[string]map[*Client]bool

	broadcast chan Message

	register chan Subscription

	unregister chan Subscription

	stats chan chan []RoomStat

	checkLimit chan LimitCheck

	store *store.Stores
}

func NewHub(st *store.Stores) *Hub {
	return &Hub{
		store:      st,
		broadcast:  make(chan Message),
		register:   make(chan Subscription),
		unregister: make(chan Subscription),
		rooms:      make(map[string]map[*Client]bool),
		stats:      make(chan chan []RoomStat),
		checkLimit: make(chan LimitCheck),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			client := sub.Client
			room := sub.Room
			if h.rooms[room] == nil {
				h.rooms[room] = make(map[*Client]bool)
			}

			if len(h.rooms[room]) >= 50 {
				log.Printf("Room %s full, rejecting client %s", room, client.username)
				client.conn.Close()
				continue
			}

			h.rooms[room][client] = true
			log.Printf("Client %s registered in room %s. Total clients in room: %d", client.username, room, len(h.rooms[room]))

		case sub := <-h.unregister:
			client := sub.Client
			room := sub.Room
			if clients, ok := h.rooms[room]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					log.Printf("Client %s unregistered from room %s. Total clients in room: %d", client.username, room, len(clients))
					if len(clients) == 0 {
						delete(h.rooms, room)
					}
				}
			}

		case message := <-h.broadcast:
			room := message.Room

			// Save to store
			h.store.Mu.Lock()
			if h.store.RoomMessages[room] == nil {
				h.store.RoomMessages[room] = []store.Message{}
			}
			h.store.RoomMessages[room] = append(h.store.RoomMessages[room], store.Message{
				ID:        uuid.NewString(), // Generate ID
				RoomID:    room,
				Sender:    message.Sender,
				Content:   message.Content,
				CreatedAt: time.Now().Format(time.RFC3339),
				Type:      "message",
			})
			h.store.Mu.Unlock()

			if clients, ok := h.rooms[room]; ok {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}

		case replyChan := <-h.stats:
			stats := make([]RoomStat, 0, len(h.rooms))
			for name, clients := range h.rooms {
				stats = append(stats, RoomStat{Name: name, Count: len(clients)})
			}
			replyChan <- stats

		case check := <-h.checkLimit:
			count := 0
			if clients, ok := h.rooms[check.Room]; ok {
				count = len(clients)
			}
			check.Result <- (count >= 50)
		}
	}
}

func (h *Hub) GetStats() []RoomStat {
	replyChan := make(chan []RoomStat)
	h.stats <- replyChan
	return <-replyChan
}

func (h *Hub) IsRoomFull(room string) bool {
	resChan := make(chan bool)
	h.checkLimit <- LimitCheck{Room: room, Result: resChan}
	return <-resChan
}
