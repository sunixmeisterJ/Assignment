package handler

import (
	"net/http/httptest"
	"rental/internal/store"
	"strings"
	"testing"
)

func TestAPIValidation(t *testing.T) {
	a := New(store.New())
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{{"GET", "/api/state", "", 200}, {"POST", "/api/people", `{"name":"Test","documentType":"Pass","documentNumber":"X"}`, 200}, {"POST", "/api/items", `{"name":"","hourlyCents":0}`, 400}, {"POST", "/api/items", `{"unknown":1}`, 400}, {"POST", "/api/people", `{} {}`, 400}, {"POST", "/api/rentals/abc/return", `{}`, 400}, {"DELETE", "/api/state", "", 405}} {
		w := httptest.NewRecorder()
		a.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body.String())
		}
	}
}
