package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/backend/structs"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn
	// manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the WebSocket
	egress chan Event
	// Client info
	clientInfo structs.User
	// Chat list order array
	ChatListOrderArray []int
	// Current active groupID //0 if no active group
	CurrentActiveGroup int
	// Current active chat recipientID //0 if no active chat recipientID
	CurrentActiveChatRecipientID int
	// List of all current joined groupIDs
	JoinedGroupIDArray []int
}

func NewClient(conn *websocket.Conn, manager *Manager, r *http.Request) *Client {
	return &Client{
		connection:                   conn,
		manager:                      manager,
		egress:                       make(chan Event),
		clientInfo:                   AddUserInfo(r),
		ChatListOrderArray:           GetChatListOrder(r),
		JoinedGroupIDArray:           GetJoinedGroupIDs(r),
		CurrentActiveGroup:           0,
		CurrentActiveChatRecipientID: 0,
	}
}

func GetJoinedGroupIDs(r *http.Request) []int {
	user, err := structs.GetUser(structs.GetUserIDfromSessions(r))
	if err != nil {
		log.Println(err)
		return []int{}
	}

	arr, err := structs.GetJoinedGroupIDArray(user.UserID)
	if err != nil {
		log.Println(err)
		return []int{}
	}
	return arr
}

func GetChatListOrder(r *http.Request) []int {
	user, err := structs.GetUser(structs.GetUserIDfromSessions(r))
	if err != nil {
		log.Println(err)
		return []int{}
	}
	arr, err := structs.GetChatListOrderArray(user.UserID)
	if err != nil {
		log.Println(err)
		return []int{}
	}
	var intArr []int
	err = json.Unmarshal(arr, &intArr)
	if err != nil {
		log.Println(err)
		return []int{}
	}
	return intArr
}

func AddUserInfo(r *http.Request) structs.User {
	user, err := structs.GetUser(structs.GetUserIDfromSessions(r))
	if err != nil {
		log.Println(err)
	}
	userInfo := structs.User{
		UserID:   user.UserID,
		Nickname: user.Nickname,
	}
	return userInfo
}

func UpdateCurrentActiveGroup(event Event, c *Client) error {
	var newCurrentActiveGroup int
	if err := json.Unmarshal(event.Payload, &newCurrentActiveGroup); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	c.CurrentActiveGroup = newCurrentActiveGroup
	return nil
}

func UpdateCurrentActiveChatRecipientID(event Event, c *Client) error {
	var newCurrentActiveChat int
	if err := json.Unmarshal(event.Payload, &newCurrentActiveChat); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	c.CurrentActiveChatRecipientID = newCurrentActiveChat
	return nil
}
