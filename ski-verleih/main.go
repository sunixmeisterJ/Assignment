package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"rental/internal/handler"
	"rental/internal/store"
	"time"
)

//go:embed web/*
var assets embed.FS

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	app := handler.New(store.New())
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", app.ServeHTTP)
	mux.Handle("/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusTemporaryRedirect)
	})
	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("Verleih läuft auf http://%s", addr)
	log.Fatal(server.ListenAndServe())
}
