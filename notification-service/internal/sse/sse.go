package sse

import (
	"errors"
	"sync"
)

type Broker struct {
	clients   map[string]chan string
	mutex     sync.Mutex
	addCh     chan Client
	removeCh  chan Client
	messageCh chan TargetedMessage
}

type Client struct {
	ID   string
	Chan chan string
}

type TargetedMessage struct {
	ID      string
	Message string
}

func NewBroker() *Broker {
	return &Broker{
		clients:   make(map[string]chan string),
		addCh:     make(chan Client),
		removeCh:  make(chan Client),
		messageCh: make(chan TargetedMessage),
	}
}

func (b *Broker) Start() {
	for {
		select {
		case client := <-b.addCh:
			b.mutex.Lock()
			b.clients[client.ID] = client.Chan
			b.mutex.Unlock()
		case client := <-b.removeCh:
			b.mutex.Lock()
			delete(b.clients, client.ID)
			close(client.Chan)
			b.mutex.Unlock()
		case msg := <-b.messageCh:
			b.mutex.Lock()
			if clientChan, exists := b.clients[msg.ID]; exists {
				clientChan <- msg.Message
			}
			b.mutex.Unlock()
		}
	}
}

func (b *Broker) AddClient(clientID string, clientChan chan string) {
	b.addCh <- Client{ID: clientID, Chan: clientChan}
}

func (b *Broker) RemoveClient(clientID string) {
	b.removeCh <- Client{ID: clientID}
}

func (b *Broker) SendMessageToClient(clientID string, msg string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if clientChan, exists := b.clients[clientID]; exists {
		clientChan <- msg
		return nil
	}

	return errors.New("client not found")
}
