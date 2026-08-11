// Command web runs the Go web examples: JSON encode/decode, HTTP client calls,
// the GitHub API, and two small HTTP servers. Pick one with -demo, e.g.:
//
//	go run ./cmd/web -demo json
//	go run ./cmd/web -demo httpd      # starts a server on :8080
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"json":       jsonDemo,
	"http-get":   httpGet,
	"github-api": githubAPI,
	"httpd":      httpd,
	"kv-store":   kvStore,
}

const depositJSON = `{"user": "Scrooge McDuck", "type": "deposit", "amount": 1000000.3}`

// Request is a bank transaction.
type Request struct {
	Login  string  `json:"user"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

func jsonDemo() {
	req := &Request{}
	if err := json.NewDecoder(bytes.NewBufferString(depositJSON)).Decode(req); err != nil {
		log.Fatalf("error: can't decode - %s", err)
	}
	fmt.Printf("got: %+v\n", req)

	resp := map[string]any{"ok": true, "balance": 8500000.0 + req.Amount}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Fatalf("error: can't encode - %s", err)
	}
}

// Job is a job description sent to httpbin.
type Job struct {
	User   string `json:"user"`
	Action string `json:"action"`
	Count  int    `json:"count"`
}

func httpGet() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalf("error: can't call httpbin.org")
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(os.Stdout, resp.Body)

	fmt.Println("----")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&Job{User: "Saitama", Action: "punch", Count: 1}); err != nil {
		log.Fatalf("error: can't encode job - %s", err)
	}
	resp, err = http.Post("https://httpbin.org/post", "application/json", &buf)
	if err != nil {
		log.Fatalf("error: can't call httpbin.org")
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(os.Stdout, resp.Body)
}

// User is GitHub user information.
type User struct {
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
}

func userInfo(login string) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s", login))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func githubAPI() {
	user, err := userInfo("tebeka")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("%+v\n", user)
}

// MathRequest is a request for a math operation.
type MathRequest struct {
	Op    string  `json:"op"`
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}

// MathResponse is the result of a MathRequest.
type MathResponse struct {
	Error  string  `json:"error"`
	Result float64 `json:"result"`
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello Gophers!")
}

func mathHandler(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	req := &MathRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &MathResponse{}
	switch req.Op {
	case "+":
		resp.Result = req.Left + req.Right
	case "-":
		resp.Result = req.Left - req.Right
	case "*":
		resp.Result = req.Left * req.Right
	case "/":
		if req.Right == 0.0 {
			resp.Error = "division by 0"
		} else {
			resp.Result = req.Left / req.Right
		}
	default:
		resp.Error = fmt.Sprintf("unknown operation: %s", req.Op)
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("can't encode %v - %s", resp, err)
	}
}

func httpd() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/math", mathHandler)
	log.Print("listening on :8080 (/hello, /math)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// Entry is a key/value store entry used for both requests and responses.
type Entry struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func kvStore() {
	var (
		db     = map[string]any{}
		dbLock sync.Mutex
	)

	sendResponse := func(entry *Entry, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entry); err != nil {
			log.Printf("error encoding %+v - %s", entry, err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		entry := &Entry{}
		if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dbLock.Lock()
		defer dbLock.Unlock()
		db[entry.Key] = entry.Value
		sendResponse(entry, w)
	})
	mux.HandleFunc("/db/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/db/"):]
		dbLock.Lock()
		defer dbLock.Unlock()
		value, ok := db[key]
		if !ok {
			http.Error(w, fmt.Sprintf("Key %q not found", key), http.StatusNotFound)
			return
		}
		sendResponse(&Entry{Key: key, Value: value}, w)
	})

	log.Print("listening on :8080 (POST /db, GET /db/<key>)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
