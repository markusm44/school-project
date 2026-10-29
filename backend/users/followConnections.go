package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "social-network/backend/database"
	"social-network/backend/structs"
	"strconv"
)

func UserStatusUpdate(urlID int, userStatus bool, r *http.Request) error {
	currentUserID := structs.GetUserIDfromSessions(r)
	if currentUserID != urlID {
		UpdateUserFollowStatus(currentUserID, urlID, userStatus)
	} else {
		UserPublicStatusUpdate(currentUserID, userStatus)
	}

	return nil
}

func UpdateUserFollowStatus(currentUserID, followedUserID int, followStatus bool) error {
	user, err := structs.GetUser(followedUserID)
	if err != nil {
		fmt.Println("ERR in GetUser: ", err)
		return err
	}

	if followStatus && user.IsProfilePublic {
		err := InsertIntoSenderDB(currentUserID, followedUserID)
		if err != nil {
			fmt.Println("ERR in InsertIntoCurrentUserDB: ", err)
			return err
		}
		err = InsertIntoRecieverDB(currentUserID, followedUserID)
		if err != nil {
			fmt.Println("ERR in InsertIntoFollowedUserDB: ", err)
			return err
		}
	} else if followStatus {
		err := SendFollowRequest(currentUserID, followedUserID, followStatus, GlobalManager)
		if err != nil {
			fmt.Println("ERR in SendFollowRequest: ", err)
			return err
		}
	} else {
		err := structs.RemoveFromFollowLists(currentUserID, followedUserID)
		if err != nil {
			fmt.Println("ERR in RemoveFromFollowedList: ", err)
			return err
		}
	}

	return nil
}

func UserPublicStatusUpdate(userID int, status bool) error {
	_, err := database.DBase.Exec("UPDATE users SET isProfilePublic = ? WHERE UserID = ? ", status, userID)
	if err != nil {
		log.Println("userStatusErr: ", err)
		return err
	}
	return nil
}

func GetProfileStatus(userID int) (bool, error) {
	var isPublic bool

	err := database.DBase.QueryRow("SELECT isProfilePublic FROM users WHERE userID = ?", userID).Scan(&isPublic)
	if err != nil {
		fmt.Println("getFollowStatusErr: ", err)
		return false, err
	}

	return isPublic, nil
}

func InsertIntoRecieverDB(senderID int, recieverID int) error {
	// Get the current followers list of the followed user
	followersList, _, err := structs.GetFollowLists(recieverID)
	if err != nil {
		return fmt.Errorf("GetFollowersLists error: %d", err)
	}

	// Append the senderID to the followers list
	followersList = append(followersList, senderID)

	// Marshal the followers list to JSON
	followersJSON, err := json.Marshal(followersList)
	if err != nil {
		return fmt.Errorf("JSON Marshal error: %d", err)
	}

	// Execute the database update
	_, err = database.DBase.Exec("UPDATE followers SET followerIDs = ? WHERE userID = ?", followersJSON, recieverID)
	if err != nil {
		return fmt.Errorf("Database update error: %d", err)
	}

	return nil
}

func InsertIntoSenderDB(senderID int, recieverID int) error {
	// Get the current user's followed users list
	_, followedUsersList, err := structs.GetFollowLists(senderID)
	if err != nil {
		return fmt.Errorf("GetFollowersLists error CU: %d", err)
	}

	// Append the recieverID to the followed users list
	followedUsersList = append(followedUsersList, recieverID)

	// Marshal the followed users list to JSON
	followedUsersJSON, err := json.Marshal(followedUsersList)
	if err != nil {
		return fmt.Errorf("JSON Marshal error: %d", err)
	}

	_, err = database.DBase.Exec("UPDATE followers SET followedUsers = ? WHERE userID = ?", followedUsersJSON, senderID)
	if err != nil {
		return fmt.Errorf("Database update error CU: %d", err)
	}

	return nil
}

func StructurizeProfileLists(r *http.Request) (structs.ProfileData, error) {
	userID, atoiErr := strconv.Atoi(r.URL.Query().Get("userID"))
	if atoiErr != nil || userID == 0 {
		fmt.Println("getUserIDConvertErrGet: ", atoiErr)
		return structs.ProfileData{}, atoiErr
	}
	userData, err := structs.GetAllUsers()
	if err != nil {
		fmt.Println("getAllUsersErr: ", err)
		return structs.ProfileData{}, err
	}

	// Create a map for faster lookup
	userMap := make(map[int]structs.User)
	for _, user := range userData {
		userMap[user.UserID] = user
	}

	var profile structs.ProfileData

	user, ok := userMap[userID]
	if ok {
		profile.UserStruct = user
	}

	followersList, followedUsersList, err := structs.GetFollowLists(userID)
	if err != nil {
		fmt.Println("GetFollowersLists err : ", err)
		return structs.ProfileData{}, err
	}

	// Populate profile.Followers
	for _, followerID := range followersList {
		if user, ok := userMap[followerID]; ok {
			profile.Followers = append(profile.Followers, structs.FollowStruct{
				UserID: followerID,
				Name:   user.FirstName + " " + user.LastName,
			})
		}
	}

	// Populate profile.Followed
	for _, followedUserID := range followedUsersList {
		if user, ok := userMap[followedUserID]; ok {
			profile.Followed = append(profile.Followed, structs.FollowStruct{
				UserID: followedUserID,
				Name:   user.FirstName + " " + user.LastName,
			})
		}
	}

	recieverID := structs.GetUserIDfromSessions(r)
	profile.Posts, err = structs.GetProfilePosts(recieverID, userID)
	if err != nil {
		fmt.Println("GetAllPostsFromUserID ERR: ", err)
		return structs.ProfileData{}, err
	}
	// Get and format all Groups a user has joined
	userJoinedGroups, err := structs.GetJoinedGroupIDArray(recieverID)
	if err != nil {
		fmt.Println("error getting GetJoinedGroupIDArray: ", err)
		return structs.ProfileData{}, err
	}

	for _, groupID := range userJoinedGroups {
		groupPosts, err := structs.GetGroupPostsForProfile(userID, groupID)
		if err != nil {
			fmt.Println("error getting group posts for { PROFILE }")
			return structs.ProfileData{}, err
		}

		// Append groupPosts to profile.GroupPosts
		profile.GroupPosts = append(profile.GroupPosts, groupPosts...)
	}
	// doesn't send passwordHash to frontend
	profile.UserStruct.PasswordHash = ""
	return profile, nil
}
