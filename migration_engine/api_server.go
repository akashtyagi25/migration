package main

import (
	"fmt"
	"net/http"
	"os"
)

// Global SSE channel to broadcast logs to the dashboard
var logChannel = make(chan string, 100)

// BroadcastEvent sends a JSON string to the dashboard
func BroadcastEvent(event string) {
	select {
	case logChannel <- event:
	default:
		// Drop message if channel is full
	}
}

// StartAPIServer launches the Web Dashboard
func StartAPIServer(port string) {
	http.Handle("/", http.FileServer(http.Dir("./web_dashboard")))

	http.HandleFunc("/api/migrate", handleMigrationStream)

	fmt.Printf("[Dashboard]: Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("[FATAL]: Server failed: %v\n", err)
		os.Exit(1)
	}
}

func handleMigrationStream(w http.ResponseWriter, r *http.Request) {
	targetLang := r.URL.Query().Get("target")
	instructions := r.URL.Query().Get("instructions")
	
	if targetLang == "" {
		targetLang = "go"
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Run migration in background
	go RunBatchMigrationUI(targetLang, instructions)

	// Stream logs
	for {
		select {
		case msg := <-logChannel:
			if msg == "DONE" {
				fmt.Fprintf(w, "event: done\ndata: {}\n\n")
				flusher.Flush()
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
