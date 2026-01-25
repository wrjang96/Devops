package users

import (
	"net/http"
	"time"

	"go-chat-mem/internal/auth"
	"go-chat-mem/internal/config"
	"go-chat-mem/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Cfg   config.Config
	Store *store.Stores
}

type joinReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h Handler) Join(c *gin.Context) {
	var req joinReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	h.Store.Mu.Lock()
	defer h.Store.Mu.Unlock()

	if _, exists := h.Store.UsersByEmail[req.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	u := store.User{ID: uuid.NewString(), Email: req.Email, PasswordHash: string(hash)}
	h.Store.UsersByEmail[u.Email] = u
	h.Store.UsersByID[u.ID] = u
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

func (h Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	h.Store.Mu.RLock()
	u, ok := h.Store.UsersByEmail[req.Email]
	h.Store.Mu.RUnlock()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	access, err := auth.SignAccessToken(h.Cfg.JWTSecret, u.ID, u.Email, h.Cfg.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	refreshPlain, refreshHash, err := store.NewRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	rec := store.RefreshTokenRecord{
		UserID:    u.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(h.Cfg.RefreshTTL),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	h.Store.Mu.Lock()
	h.Store.RefreshByHash[refreshHash] = rec
	h.Store.Mu.Unlock()

	c.SetCookie(h.Cfg.RefreshCookieKey, refreshPlain, int(h.Cfg.RefreshTTL.Seconds()), "/", "", h.Cfg.CookieSecure, true)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": access,
		"tokenType":   "Bearer",
		"expiresIn":   int(h.Cfg.AccessTTL.Seconds()),
	})
}

func (h Handler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie(h.Cfg.RefreshCookieKey)
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	hash := store.HashRefreshToken(rt)

	h.Store.Mu.RLock()
	rec, ok := h.Store.RefreshByHash[hash]
	h.Store.Mu.RUnlock()
	if !ok || rec.Revoked || time.Now().After(rec.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	h.Store.Mu.RLock()
	u, ok := h.Store.UsersByID[rec.UserID]
	h.Store.Mu.RUnlock()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	access, err := auth.SignAccessToken(h.Cfg.JWTSecret, u.ID, u.Email, h.Cfg.AccessTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	// rotation
	h.Store.Mu.Lock()
	old := h.Store.RefreshByHash[hash]
	old.Revoked = true
	h.Store.RefreshByHash[hash] = old
	h.Store.Mu.Unlock()

	newPlain, newHash, err := store.NewRefreshToken()
	if err == nil {
		h.Store.Mu.Lock()
		h.Store.RefreshByHash[newHash] = store.RefreshTokenRecord{
			UserID:    u.ID,
			TokenHash: newHash,
			ExpiresAt: time.Now().Add(h.Cfg.RefreshTTL),
			Revoked:   false,
			CreatedAt: time.Now(),
		}
		h.Store.Mu.Unlock()

		c.SetCookie(h.Cfg.RefreshCookieKey, newPlain, int(h.Cfg.RefreshTTL.Seconds()), "/", "", h.Cfg.CookieSecure, true)
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": access,
		"tokenType":   "Bearer",
		"expiresIn":   int(h.Cfg.AccessTTL.Seconds()),
	})
}
