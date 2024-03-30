package concurrency_pattern

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

// --- server ----
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Handler started")
	defer log.Printf("Handler stopped")
	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintf(w, "hello")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Test13_server(t *testing.T) {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

// client

func Test13_client(t *testing.T) {
	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, "GET", "https://www.qq.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
