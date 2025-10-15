package main

import (
	"log"
	"net/http"
	"time"

	"chat-server/internal/config"
	"chat-server/internal/handlers"
	"chat-server/internal/middleware"
	"chat-server/internal/room"
	"chat-server/internal/user"
)

func main() {
	cfg := config.New()

	roomManager := room.NewManager()

	userManager := user.NewManager()

	handlers := handlers.New(roomManager, userManager, cfg)

	rateLimiter := middleware.NewRateLimiter(100, 60*time.Second)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	mux.HandleFunc("/", handlers.ServeHome)
	mux.HandleFunc("/ws", handlers.HandleWebSocket)
	mux.HandleFunc("/api/rooms", handlers.GetRooms)

	handler := middleware.CORS(mux)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.RateLimit(rateLimiter)(handler)

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
