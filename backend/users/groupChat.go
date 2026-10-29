package users

import (
	"encoding/json"
	"fmt"
	"social-network/backend/structs"
)

// SendMessageHandler will send out a message to all other participants in the group
func SendGroupMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var groupMessage SendGroupMessageEvent
	if err := json.Unmarshal(event.Payload, &groupMessage); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	err := structs.AddGroupMessages(groupMessage.SenderID, groupMessage.SenderName, groupMessage.GroupID, groupMessage.Message, groupMessage.Date)
	if err != nil {
		return err
	}

	for clientFromClientList := range c.manager.clients {
		//if the chat is active send the eventPayload through websocket
		if clientFromClientList.CurrentActiveGroup == groupMessage.GroupID {
			// Prepare an Outgoing Message to others
			var broadMessage SendGroupMessageEvent
			broadMessage.SenderName = groupMessage.SenderName
			broadMessage.Message = groupMessage.Message
			broadMessage.SenderID = groupMessage.SenderID
			broadMessage.GroupID = groupMessage.GroupID
			broadMessage.Date = groupMessage.Date

			data, err := json.Marshal(broadMessage)
			if err != nil {
				return fmt.Errorf("failed to marshal broadcast message: %v", err)
			}

			// Place payload into an Event
			var outgoingEvent Event
			outgoingEvent.Payload = data
			outgoingEvent.Type = EventNewGroupMessage

			// Broadcast to group
			clientFromClientList.egress <- outgoingEvent
		}
	}
	return nil
}
