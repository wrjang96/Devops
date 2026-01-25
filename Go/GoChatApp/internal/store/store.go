package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"sync"
	"time"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type RefreshTokenRecord struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

type Room struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	Members   int    `json:"members"`
}

type Stores struct {
	Mu sync.RWMutex

	UsersByEmail map[string]User
	UsersByID    map[string]User

	RefreshByHash map[string]RefreshTokenRecord

	Rooms       map[string]Room
	RoomMembers map[string]map[string]struct{}
}

func New() *Stores {
	return &Stores{
		UsersByEmail:  map[string]User{},
		UsersByID:     map[string]User{},
		RefreshByHash: map[string]RefreshTokenRecord{},
		Rooms:         map[string]Room{},
		RoomMembers:   map[string]map[string]struct{}{},
	}
}

func NewRefreshToken() (plain, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	plain = base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])
	return plain, hash, nil
}

func HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
