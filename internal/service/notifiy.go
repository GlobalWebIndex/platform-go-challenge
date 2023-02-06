package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	//"golang.org/x/crypto/bcrypt"
)

type NotifyService interface {
	SendMessage(data string) error
}

type notifyService struct {
	conn *websocket.Conn
}

// SendMessage implements NotifyService
func (n *notifyService) SendMessage(data string) error {
	err := ws.WriteMessage(1, []byte(data))
	if err != nil {
		log.Println(err)
	}
	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ws *websocket.Conn

func NewNotifyService() NotifyService {
	http.HandleFunc("/ws", wsEndpoint)
	return &notifyService{
		conn: ws,
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// helpful log statement to show connections
	log.Println("Client Connected")
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
