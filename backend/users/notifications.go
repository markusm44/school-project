package users

import (
	"encoding/json"
	"fmt"
	"social-network/backend/structs"
)

const (
	NotifNameForFollowRequest         = "follow_request"
	NotifNameForRequestingToJoinGroup = "request_to_join_group"
	NotifNameForInvitingToJoinGroup   = "invitation_to_join_group"
	NotifNameForGroupEvent            = "group_event_notification"
)

type GroupMessageNotification struct {
	GroupID    int    `json:"groupID"`
	GroupName  string `json:"groupName"`
	SenderID   int    `json:"senderID"`
	SenderName string `json:"senderName"`
	Date       string `json:"date"`
}

func SendNotification(c *Client, notificationData interface{}) error {
	// Marshal the notification data
	data, err := json.Marshal(notificationData)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewNotification
	// Broadcast to recipient
	c.egress <- outgoingEvent
	return nil
}

func DeleteNotificationHandler(event Event, c *Client) error {
	var notificationID int
	err := json.Unmarshal(event.Payload, &notificationID)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}
	err = structs.DeleteNotification(notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification from db: %v", err)
	}
	return nil
}

func SendFollowRequest(currentUserID, followedUserID int, startedFollowing bool, manager *Manager) error {
	user, err := structs.GetUser(currentUserID)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	notificationID, err := structs.AddNotification(NotifNameForFollowRequest, content, currentUserID, followedUserID)
	if err != nil {
		return err
	}
	notificationData := structs.NewNotificationEvent{
		Subtype: NotifNameForFollowRequest,
		Body: structs.Notification{
			NotificationID:      notificationID,
			NotificationSubtype: NotifNameForFollowRequest,
			SenderID:            currentUserID,
			RecieverID:          followedUserID,
			Content:             content,
			Created:             "",
		},
	}
	for clientFromClientList := range manager.clients {
		if clientFromClientList.clientInfo.UserID == followedUserID {
			err = SendNotification(clientFromClientList, notificationData)
			if err != nil {
				return fmt.Errorf("failed to send event: %v", err)
			}
		}
	}
	return nil
}
