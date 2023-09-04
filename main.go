package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Chat struct {
	messages []string
	mu       sync.Mutex
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Messages []string `json:"messages"`
}

func main() {
	chat := &Chat{}

	r := mux.NewRouter()
	r.HandleFunc("/send", chat.send).Methods("GET")
	r.HandleFunc("/receive", chat.receive).Methods("GET")

	corsHandler := cors.AllowAll().Handler(r)

	fmt.Println("Chat server started on :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		fmt.Println(err)
	}
}

func (c *Chat) AddMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, message)
}

func (c *Chat) GetMessages() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.messages
}

func (chat *Chat) receive(w http.ResponseWriter, r *http.Request) {
	messages := chat.GetMessages()

	response := ResponseData{Messages: messages}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (chat *Chat) send(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message != "" {
		chat.AddMessage(message)
		response := MessageResponse{Message: fmt.Sprintf("Message sent: %s", message)}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		http.Error(w, "Message parameter is required", http.StatusBadRequest)
	}
}
