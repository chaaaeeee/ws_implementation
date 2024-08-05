package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok { // is there any room with client's roomID as it key?
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok { // is there any client with the key of client's ID? id is equal to key no?
					r.Clients[cl.ID] = cl // if not then store client as one of the client that is in this room, so key is indeed equals to id
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "User left the room",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					} // broadcast a msg that the client left the room
					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}
		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomID]; ok {
				for _, cl := range h.Rooms[msg.RoomID].Clients {
					cl.Message <- msg
				}
			}
		}
	}

}
