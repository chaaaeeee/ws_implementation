package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // establish connection to websocket
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// /ws/JoinRoom/roomId?userId=1&username=user
	// take the value from url above
	roomID := c.Param("roomId")
	userID := c.Query("userId")
	username := c.Query("username")

	// store it into client struct
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10), // this a homework
		ID:       userID,
		RoomID:   roomID,
		Username: username,
	}

	// make a message struct
	msg := &Message{
		Content:  "A user has joined the room!",
		RoomID:   roomID,
		Username: username,
	}

	// register new client through the register channel
	h.hub.Register <- cl
	// broadcast the message
	h.hub.Broadcast <- msg

	// writeMessage()

	// readMessage()
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

func (h *Handler) GetRooms(c *gin.Context) {
	c.JSON(http.StatusOK, h.hub.Rooms)
}
