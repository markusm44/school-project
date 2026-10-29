package users

import (
	"encoding/json"
	"fmt"
	"log"
	"social-network/backend/structs"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	// INBOUND EVENTS

	// EventJoinGroup is the event for getting the usersID and groupID for sending group event notification.
	EventSendGroupEventNotification = "notify_of_group_event"
	// EventJoinGroup is the event for getting the usersID and groupID for inviting to Group.
	EventInviteToGroup = "invite_to_group"
	// EventJoinGroup is the event for getting the usersID and groupID for joining.
	EventJoinGroup = "join_group"
	// EventGetActiveGroupIDArray is the event for triggering sending JoinedGroupIDArray to frontend
	EventGetActiveGroupIDArray = "get_joined_groupID_array"
	// EventUpdateActiveGroupInBackend is the event for updating the CurrentActiveGroup field in the Client struct
	EventUpdateActiveGroupInBackend = "update_active_group"
	// EventUpdateActiveChatInClient is the event for updating the CurrentActiveChatRecipientID field in the Client struct
	EventUpdateActiveChatInClient = "update_active_chat"
	// EventSendMessage is the event name for new chat messages sent from the frontend
	EventSendMessage = "send_message"
	// EventSendMessage is the event name for new group messages sent
	EventSendGroupMessage = "send_group_message"
	// Event to update the allUsersFromDB field in manager
	EventUpdateAllUsersInDB = "update_all_users_from_DB"
	// Event to update the client.ChatListOrderArray value
	EventUpdateChatListOrderArray = "chat_list_order_array"
	// Event to send a join request to the group creator
	EventRequestToJoinGroup = "request_to_join_group"
	// Event to delete notification from database
	EventDeleteNotification = "delete_notification"

	// OUTBOUND EVENTS

	// EventSendMessage is the event name for all notifications
	EventNewNotification = "new_notification"
	// EventSendActiveGroupIDArray is the event for sending JoinedGroupIDArray to frontend
	EventSendActiveGroupIDArray = "send_joined_groupID_array"
	// EventNewMessage is a response to send_message
	EventNewMessage = "new_message"
	// EventNewMessage is a response to "send_group_message"
	EventNewGroupMessage = "new_group_message"
	// EventOnlineUsersList on every adding or removing of a client update the chat user list
	EventOnlineUsersList = "clients_online_status"
	// EventOnlineUsersList on every adding or removing of a client update the chat user list
	EventNewGroupEventNotification = "group_event_notification"
)

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventUpdateChatListOrderArray] = UpdateChatListOrderArray
	m.handlers[EventUpdateAllUsersInDB] = UpdateAllUsersFromDB
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventSendGroupMessage] = SendGroupMessageHandler
	m.handlers[EventUpdateActiveGroupInBackend] = UpdateCurrentActiveGroup
	m.handlers[EventUpdateActiveChatInClient] = UpdateCurrentActiveChatRecipientID
	m.handlers[EventGetActiveGroupIDArray] = GetJoinedGroupIDArray
	m.handlers[EventJoinGroup] = JoinGroup
	m.handlers[EventRequestToJoinGroup] = RequestToJoinGroup
	m.handlers[EventInviteToGroup] = InviteToJoinGroup
	m.handlers[EventSendGroupEventNotification] = SendGroupEventNotification
	m.handlers[EventDeleteNotification] = DeleteNotificationHandler
}

type SendGroupMessageEvent struct {
	GroupID    int    `json:"groupID"`
	Message    string `json:"message"`
	SenderID   int    `json:"senderID"`
	SenderName string `json:"senderName"`
	Date       string `json:"date"`
}

// SendMessageEvent is the payload sent in the send_message event
type SendMessageEvent struct {
	Message     string `json:"message"`
	SenderID    int    `json:"senderID"`
	RecipientID int    `json:"recipientID"`
	Date        string `json:"date"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent string `json:"sent"`
}

// SendMessageHandler will send out a message to recipientID in the chat
func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	err := structs.AddMessages(chatevent.SenderID, chatevent.RecipientID, chatevent.Message, chatevent.Date)
	if err != nil {
		return err
	}

	// Prepare an Outgoing Message to others
	var broadMessage NewMessageEvent
	broadMessage.Message = chatevent.Message
	broadMessage.SenderID = chatevent.SenderID
	broadMessage.RecipientID = chatevent.RecipientID
	broadMessage.Date = chatevent.Date

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage
	// Broadcast to recipient
	for clientFromClientList := range c.manager.clients {
		if clientFromClientList.clientInfo.UserID == chatevent.RecipientID {
			clientFromClientList.egress <- outgoingEvent
		}
	}
	return nil
}

func updateChatList(c *Client) error {
	//Get all registerd users from db
	allUsersFromDB := c.manager.allUsersFromDB
	// make a slice of online users
	var onlineClientsSlice []structs.User
	for client := range c.manager.clients {
		onlineClientsSlice = append(onlineClientsSlice, client.clientInfo)
	}

	type everyUserStatus struct {
		structs.User
		Online bool `json:"onlineStatus"`
	}

	var clientOnlineStatusData []everyUserStatus
	online := false

	for _, user := range allUsersFromDB {
		var onlineClient everyUserStatus
		onlineClient.UserID = user.UserID
		onlineClient.FirstName = user.FirstName
		onlineClient.LastName = user.LastName
		for _, onlineUser := range onlineClientsSlice {
			if user.UserID == onlineUser.UserID {
				online = true
				break
			}
		}
		if online {
			onlineClient.Online = true
		} else {
			onlineClient.Online = false
		}
		clientOnlineStatusData = append(clientOnlineStatusData, onlineClient)
		online = false
	}
	payload, err := json.Marshal(clientOnlineStatusData)
	if err != nil {
		log.Println(err)
		return err
	}

	var outgoingEvent Event
	// give Event type
	outgoingEvent.Type = EventOnlineUsersList
	outgoingEvent.Payload = payload

	// Write a message to all online client connections
	for clientFromClientList := range c.manager.clients {
		clientFromClientList.egress <- outgoingEvent
	}
	return nil
}

func UpdateChatListOrderArray(event Event, c *Client) error {
	var chatListOrderArray []int
	if err := json.Unmarshal(event.Payload, &chatListOrderArray); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	c.ChatListOrderArray = chatListOrderArray
	return nil
}
