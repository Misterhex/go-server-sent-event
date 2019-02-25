package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	http.Handle("/events/", http.HandlerFunc(handler))

	// Start the server and listen forever on port 8000.
	if err := http.ListenAndServe(":10090", nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	notify := w.(http.CloseNotifier).CloseNotify()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

OuterLoop:
	for {
		select {
		case <-time.After(2 * time.Second):
			msg := fmt.Sprintf("the time is %s", time.Now().String())
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			f.Flush()
		case <-notify:
			log.Println("HTTP connection just closed.")
			break OuterLoop
		}
	}

	fmt.Println("exiting goroutine ...")
}
