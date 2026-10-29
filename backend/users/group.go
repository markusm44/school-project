package users

import (
	"encoding/json"
	"fmt"
	"social-network/backend/structs"
	"strconv"
	"strings"
)

func updateGroupMemberArray(groupID, userID int) error {
	group, err := structs.GetGroup(groupID)
	if err != nil {
		return err
	}

	var members []int
	if err := json.Unmarshal(group.Members, &members); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	//check if the user already is a member
	for _, member := range members {
		if member == userID {
			fmt.Printf("User %v has already joined group with ID %v\n", userID, groupID)
			return nil
		}
	}
	members = append(members, userID)

	err = structs.UpdateGroupMembers(members, groupID)
	if err != nil {
		return err
	}
	return nil
}

func updateUserJoinedGroupArray(groupID, userID int) error {
	joinedGroupIDArray, err := structs.GetJoinedGroupIDArray(userID)
	if err != nil {
		return err
	}
	//check if the user already is a member
	for _, groupIDFromList := range joinedGroupIDArray {
		if groupIDFromList == groupID {
			fmt.Printf("User %v has already joined group with ID %v\n", userID, groupID)
			return nil
		}
	}
	joinedGroupIDArray = append(joinedGroupIDArray, groupID)

	err = structs.UpadateJoinedGroupIDArray(joinedGroupIDArray, userID)
	if err != nil {
		return err
	}

	return nil
}

func JoinGroup(event Event, c *Client) error {
	var payload string
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	arr := strings.Split(payload, "|")
	groupID, err := strconv.Atoi(arr[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	userID, err := strconv.Atoi(arr[1])
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Updating the groups member array
	err = updateGroupMemberArray(groupID, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// update the users joinedGroupIDArray
	err = updateUserJoinedGroupArray(groupID, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = SendJoinedGroupIDArray(c, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetJoinedGroupIDArray(event Event, c *Client) error {
	var userID int
	if err := json.Unmarshal(event.Payload, &userID); err != nil {
		return fmt.Errorf("bad payload in SendJoinedGroupIDArray() request: %v", err)
	}
	err := SendJoinedGroupIDArray(c, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SendJoinedGroupIDArray(c *Client, userID int) error {
	joinedGroupIDArray, err := structs.GetJoinedGroupIDArray(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp, err := json.Marshal(joinedGroupIDArray)
	if err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = resp
	outgoingEvent.Type = EventSendActiveGroupIDArray
	// Broadcast to recipient
	for clientFromClientList := range c.manager.clients {
		// If the client is onine update their joinedGroupIDArray
		if clientFromClientList.clientInfo.UserID == userID {
			clientFromClientList.egress <- outgoingEvent
			break
		}
	}
	return nil
}

func SendGroupEventNotification(event Event, c *Client) error {
	var groupID int
	if err := json.Unmarshal(event.Payload, &groupID); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	group, err := structs.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	var groupMembers []int

	err = json.Unmarshal(group.Members, &groupMembers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	content := fmt.Sprintf("%s|%v", group.Name, groupID)

	//add and send notificatons
	for _, groupMemberID := range groupMembers {
		notificationID, err := structs.AddNotification(EventNewGroupEventNotification, content, group.CreatorID, groupMemberID)
		if err != nil {
			return fmt.Errorf("failed to add notification: %v", err)
		}
		for clientFromClientList := range c.manager.clients {
			if clientFromClientList.clientInfo.UserID == groupMemberID && c.clientInfo.UserID != groupMemberID {
				notificationData := structs.NewNotificationEvent{
					Subtype: EventNewGroupEventNotification,
					Body: structs.Notification{
						NotificationID:      notificationID,
						NotificationSubtype: EventNewGroupEventNotification,
						SenderID:            group.CreatorID,
						RecieverID:          groupMemberID,
						Content:             content,
						Created:             "",
					},
				}
				err = SendNotification(clientFromClientList, notificationData)
				if err != nil {
					return fmt.Errorf("failed to send event: %v", err)
				}
				break
			}
		}
	}
	return nil
}

func InviteToJoinGroup(event Event, c *Client) error {
	var payload string
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	arr := strings.Split(payload, "|")
	groupID, err := strconv.Atoi(arr[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	userID, err := strconv.Atoi(arr[1])
	if err != nil {
		fmt.Println(err)
		return err
	}

	group, err := structs.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	content := fmt.Sprintf("%s|%v", group.Name, groupID)

	notificationID, err := structs.AddNotification(NotifNameForInvitingToJoinGroup, content, c.clientInfo.UserID, userID)
	if err != nil {
		return fmt.Errorf("failed to add notification: %v", err)
	}

	notificationData := structs.NewNotificationEvent{
		Subtype: NotifNameForInvitingToJoinGroup,
		Body: structs.Notification{
			NotificationID:      notificationID,
			NotificationSubtype: NotifNameForInvitingToJoinGroup,
			SenderID:            c.clientInfo.UserID,
			RecieverID:          userID,
			Content:             content,
			Created:             "",
		},
	}

	for clientFromClientList := range c.manager.clients {
		if clientFromClientList.clientInfo.UserID == userID {
			err = SendNotification(clientFromClientList, notificationData)
			if err != nil {
				return fmt.Errorf("failed to send event: %v", err)
			}
		}
	}
	return nil
}

func RequestToJoinGroup(event Event, c *Client) error {
	var payload string
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	arr := strings.Split(payload, "|")
	groupID, err := strconv.Atoi(arr[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	userID, err := strconv.Atoi(arr[1])
	if err != nil {
		fmt.Println(err)
		return err
	}
	group, err := structs.GetGroup(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %v", err)
	}

	user, err := structs.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	content := fmt.Sprintf("%s|%s|%v|%v", user.FirstName+" "+user.LastName, group.Name, groupID, userID)

	notificationID, err := structs.AddNotification(NotifNameForRequestingToJoinGroup, content, userID, group.CreatorID)
	if err != nil {
		return fmt.Errorf("failed to add notification: %v", err)
	}

	notificationData := structs.NewNotificationEvent{
		Subtype: NotifNameForRequestingToJoinGroup,
		Body: structs.Notification{
			NotificationID:      notificationID,
			NotificationSubtype: NotifNameForRequestingToJoinGroup,
			SenderID:            userID,
			RecieverID:          group.CreatorID,
			Content:             content,
			Created:             "",
		},
	}
	for clientFromClientList := range c.manager.clients {
		if clientFromClientList.clientInfo.UserID == group.CreatorID {
			err = SendNotification(clientFromClientList, notificationData)
			if err != nil {
				return fmt.Errorf("failed to send event: %v", err)
			}
		}
	}
	return nil
}
