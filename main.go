package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Chat struct {
	messages []string
	mu       sync.Mutex
}

func main() {
	chat := &Chat{}

	http.HandleFunc("/send", chat.send)

	http.HandleFunc("/receive", chat.receive)

	fmt.Println("Chat server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

func (c *Chat) AddMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, message)
}
// test
func (c *Chat) GetMessages() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.messages
}

func (chat *Chat) receive(w http.ResponseWriter, r *http.Request) {
	messages := chat.GetMessages()
	for _, message := range messages {
		fmt.Fprintf(w, "%s\n", message)
	}
}

func (chat *Chat) send(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message != "" {
		chat.AddMessage(message)
		fmt.Fprintf(w, "Message sent: %s\n", message)
	} else {
		http.Error(w, "Message parameter is required", http.StatusBadRequest)
	}
}
