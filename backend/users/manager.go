package users

import (
	"errors"
	"log"
	"net/http"
	"social-network/backend/structs"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

var (
	//websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func checkOrigin(r *http.Request) bool {
	// make sure that the origin of connection is from a trusted source
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8081":
		return true
	default:
		return false
	}
}

type Manager struct {
	clients ClientList
	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
	// All userID and Usernames currnetly in the DB
	allUsersFromDB []structs.User
}

var GlobalManager *Manager

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	val, err := structs.GetAllUsers()
	if err != nil {
		log.Println(err)
	}
	m := &Manager{
		clients:        make(ClientList),
		handlers:       make(map[string]EventHandler),
		allUsersFromDB: val,
	}
	m.setupEventHandlers()
	GlobalManager = m
	return m
}

func UpdateAllUsersFromDB(event Event, c *Client) error {
	// Lock so we can manipulate
	c.manager.Lock()
	defer c.manager.Unlock()

	// Update value
	val, err := structs.GetAllUsers()
	if err != nil {
		log.Println("ERROR in UpdateAllUsersFromDB ", err)
		return err
	}
	c.manager.allUsersFromDB = val
	updateChatList(c)
	return nil
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	// upgrade this connection to a WebSocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m, r)
	// Add the newly created client to the manager
	m.addClient(client)
	// helpful log statement to show connections
	log.Println("Client Connected")

	go client.readMessages()
	go client.writeMessages()

	//update the ChatList of evey connection with new users or new online statuses
	err = updateChatList(client)
	if err != nil {
		log.Println("ERROR in onlineClientsHandler", err)
	}
}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	defer m.Unlock()
	// Lock so we can manipulate
	m.Lock()
	// Add Client
	m.clients[client] = true

}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	err := updateChatList(client)
	m.Lock()
	defer m.Unlock()
	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
	} else {
		// if client has been deleted return to avoid running updateChatList
		return
	}
	if err != nil {
		log.Println("ERROR in removeClient()", err)
	}
	// Write the ChatListOrderArray to DB when user disconnects/logs out
	err = structs.AddChatListOrderArray(client.ChatListOrderArray, client.clientInfo.UserID)
	if err != nil {
		log.Println("ERROR in removeClient()", err)
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}
