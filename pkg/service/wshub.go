package service

type HubSet struct {
	// Registered clients.
	clients map[*WsSignalConnClient]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsSignalConnClient

	// Unregister requests from clients.
	unregister chan *WsSignalConnClient
}

func newHubSet() *HubSet {
	return &HubSet{
		broadcast:  make(chan []byte),
		register:   make(chan *WsSignalConnClient),
		unregister: make(chan *WsSignalConnClient),
		clients:    make(map[*WsSignalConnClient]bool),
	}
}

func (h *HubSet) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.ResponseBuffer)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.ResponseBuffer <- message:
				default:
					close(client.ResponseBuffer)
					delete(h.clients, client)
				}
			}
		}
	}
}
