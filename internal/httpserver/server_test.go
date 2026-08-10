package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	rootHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got, want := rec.Body.String(), "Hello world!"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
	if ct := res.Header.Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/plain; charset=utf-8", ct)
	}
}

func TestBooksHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()

	booksHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}

	var got BooksResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := BooksResponse{Total: 0, Count: 0, Books: []Book{}}
	if got.Total != want.Total || got.Count != want.Count {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if got.Books == nil {
		t.Error("books should be an empty slice, not null")
	}
}

// TestRoutes exercises the wired-up mux, including method matching and 404s.
func TestRoutes(t *testing.T) {
	srv := New(Config{Addr: ":0"})

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"root ok", http.MethodGet, "/", http.StatusOK},
		{"books ok", http.MethodGet, "/books", http.StatusOK},
		{"post books not allowed", http.MethodPost, "/books", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/nope", http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			srv.Handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("%s %s = %d, want %d", tc.method, tc.path, rec.Code, tc.wantStatus)
			}
		})
	}
}
