package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notification-service/internal/sse"
)

var Broker *sse.Broker

func InitSSEHandler() {
	Broker = sse.NewBroker()
	go Broker.Start()
}

func SSEHandler(c *gin.Context) {
	clientID := c.Param("id")

	// Create a new channel for this client
	clientChan := make(chan string)

	// Add the client to the broker
	Broker.AddClient(clientID, clientChan)

	// Add CORS headers
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Update with your frontend URL
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Close the connection when done
	defer func() {
		Broker.RemoveClient(clientID)
		close(clientChan)
	}()

	// Listen for messages and send them to the client
	for {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return // Client channel closed
			}
			_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
			if err != nil {
				return // Client connection closed
			}
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return // Client connection closed by the client
		}
	}
}
