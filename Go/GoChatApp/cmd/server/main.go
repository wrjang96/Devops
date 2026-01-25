package main

import (
	"go-chat-mem/internal/config"
	httpx "go-chat-mem/internal/http"
	"go-chat-mem/internal/rooms"
	"go-chat-mem/internal/store"
	"go-chat-mem/internal/users"
	"go-chat-mem/internal/ws"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad() // config를 로드한다.
	st := store.New()

	uh := users.Handler{Cfg: cfg, Store: st}
	rh := rooms.Handler{Store: st}

	hub := ws.NewHub() // 웹소켓허브를생성한다.
	go hub.Run()

	r := httpx.NewRouter(cfg, uh, rh, hub) // 라우터를 생성하는 함수

	log.Println("listening on :" + cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
