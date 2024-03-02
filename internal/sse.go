package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendEvents(c chan CPacket, mc chan CHost) {
	http.HandleFunc("/packets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

		for packet := range c {
			data, err := json.Marshal(packet)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Send the event
			_, _ = w.Write([]byte("data: " + string(data) + "\n\n"))
			w.(http.Flusher).Flush()
		}
	})

	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

		for host := range mc {
			data, err := json.Marshal(host)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, _ = w.Write([]byte("data: " + string(data) + "\n\n"))
			w.(http.Flusher).Flush()
		}
	})

	fmt.Println("SSE on PORT 4444")
	http.ListenAndServe(":4444", nil)
}
