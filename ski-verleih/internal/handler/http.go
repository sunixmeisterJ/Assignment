package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rental/internal/domain"
	"rental/internal/store"
	"strconv"
	"strings"
)

type App struct{ store *store.Store }

func New(s *store.Store) *App { return &App{store: s} }
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method == "GET" && r.URL.Path == "/api/state" {
		a.reply(w, a.store.Snapshot(), nil)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Endpunkt oder Methode ungültig"}); err != nil {
			log.Printf("API-Antwort schreiben: %v", err)
		}
		return
	}
	if origin := r.Header.Get("Origin"); origin != "" && origin != "http://"+r.Host && origin != "https://"+r.Host {
		http.Error(w, `{"error":"Fremder Ursprung"}`, http.StatusForbidden)
		return
	}
	switch r.URL.Path {
	case "/api/items":
		var v domain.Item
		if !a.decode(w, r, &v) {
			return
		}
		result, err := a.store.SaveItem(v)
		a.reply(w, result, err)
	case "/api/people":
		var v domain.Person
		if !a.decode(w, r, &v) {
			return
		}
		result, err := a.store.AddPerson(v)
		a.reply(w, result, err)
	case "/api/rentals":
		var v domain.Rental
		if !a.decode(w, r, &v) {
			return
		}
		result, err := a.store.Reserve(v)
		a.reply(w, result, err)
	default:
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 4 || parts[0] != "api" || parts[1] != "rentals" {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			a.reply(w, nil, fmt.Errorf("Ungültige ID"))
			return
		}
		result, err := a.store.Transition(id, parts[3])
		a.reply(w, result, err)
	}
}
func (a *App) decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		a.reply(w, nil, fmt.Errorf("Ungültige Eingabe: %w", err))
		return false
	}
	if err := d.Decode(&struct{}{}); err != io.EOF {
		a.reply(w, nil, fmt.Errorf("Genau ein JSON-Objekt erwartet"))
		return false
	}
	return true
}
func (a *App) reply(w http.ResponseWriter, v any, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		v = map[string]string{"error": err.Error()}
	}
	if writeErr := json.NewEncoder(w).Encode(v); writeErr != nil {
		log.Printf("API-Antwort schreiben: %v", writeErr)
	}
}
