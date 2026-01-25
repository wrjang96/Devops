package rooms

import (
	"net/http"
	"time"

	"go-chat-mem/internal/auth"
	"go-chat-mem/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Store *store.Stores
}

type createRoomReq struct {
	Name string `json:"name"`
}

func (h Handler) List(c *gin.Context) {
	_ = auth.MustGetUser(c)

	h.Store.Mu.RLock()
	defer h.Store.Mu.RUnlock()

	out := make([]store.Room, 0, len(h.Store.Rooms))
	for _, rm := range h.Store.Rooms {
		rm.Members = len(h.Store.RoomMembers[rm.ID])
		out = append(out, rm)
	}
	c.JSON(http.StatusOK, out)
}

func (h Handler) Create(c *gin.Context) {
	u := auth.MustGetUser(c)

	var req createRoomReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	rm := store.Room{
		ID:        uuid.NewString(),
		Name:      req.Name,
		CreatedBy: u.UserID,
		CreatedAt: time.Now().Format(time.RFC3339),
		Members:   0,
	}

	h.Store.Mu.Lock()
	h.Store.Rooms[rm.ID] = rm
	if h.Store.RoomMembers[rm.ID] == nil {
		h.Store.RoomMembers[rm.ID] = map[string]struct{}{}
	}
	h.Store.Mu.Unlock()

	c.JSON(http.StatusCreated, rm)
}

func (h Handler) Join(c *gin.Context) {
	u := auth.MustGetUser(c)
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room id"})
		return
	}

	h.Store.Mu.Lock()
	defer h.Store.Mu.Unlock()

	rm, ok := h.Store.Rooms[roomID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	if h.Store.RoomMembers[roomID] == nil {
		h.Store.RoomMembers[roomID] = map[string]struct{}{}
	}
	h.Store.RoomMembers[roomID][u.UserID] = struct{}{}
	rm.Members = len(h.Store.RoomMembers[roomID])
	h.Store.Rooms[roomID] = rm

	c.JSON(http.StatusOK, gin.H{"roomId": roomID, "members": rm.Members})
}

func (h Handler) Leave(c *gin.Context) {
	u := auth.MustGetUser(c)
	roomID := c.Param("id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room id"})
		return
	}

	h.Store.Mu.Lock()
	defer h.Store.Mu.Unlock() // defer 키워드는 함수가 종료되기 전에 실행되는 코드를 지정하는 데 사용

	rm, ok := h.Store.Rooms[roomID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	if h.Store.RoomMembers[roomID] != nil {
		delete(h.Store.RoomMembers[roomID], u.UserID)
	}
	rm.Members = len(h.Store.RoomMembers[roomID])
	h.Store.Rooms[roomID] = rm

	c.JSON(http.StatusOK, gin.H{"roomId": roomID, "members": rm.Members})
}
