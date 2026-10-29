package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	database "social-network/backend/database"
	"social-network/backend/sessions"
	"social-network/backend/structs"
	"social-network/backend/users"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func profileDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		http.Error(w, "404 method not allowed", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		changeImage := r.URL.Query().Get("changeImage")
		if changeImage == "true" {
			// fmt.Println("*******profileHandler CHANGE IMAGE*********")
			var imageInfo structs.Image
			err := json.NewDecoder(r.Body).Decode(&imageInfo)
			if err != nil {
				fmt.Println("DecoderErr: ", err)
				http.Error(w, "400 bad request.", http.StatusBadRequest)
				return
			}

			userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
			if err != nil {
				fmt.Println(err)
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}

			image, err := structs.ImageNameHandler(imageInfo.Name, imageInfo.Base64String, "./media/userImages/")
			if err != nil {
				fmt.Println("Error processing image:", err)
				http.Error(w, "Error processing image.", http.StatusInternalServerError)
				return
			}

			_, err = database.DBase.Exec("UPDATE users SET image = ? WHERE userID = ?", image, userID)
			if err != nil {
				fmt.Println("Error updating user image:", err)
				http.Error(w, "Error updating user image.", http.StatusInternalServerError)
				return
			}
			// Modify to return a JSON object
			resp, err := json.Marshal(image)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else {
			// fmt.Println("*******profileHandler running Post*********")

			userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
			if err != nil {
				fmt.Println(err)
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}

			var userPublicStatus structs.ProfileInfo

			err = json.NewDecoder(r.Body).Decode(&userPublicStatus)
			if err != nil {
				fmt.Println("DecoderErr: ", err)
				http.Error(w, "400 bad request.", http.StatusBadRequest)
				return
			}

			err = users.UserStatusUpdate(userID, userPublicStatus.PublicFollowBool, r)
			if err != nil {
				log.Println("userUpdateError: ", err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}

			resp, err := json.Marshal(userPublicStatus.PublicFollowBool)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	} else if r.Method == "GET" {
		// fmt.Println("*******profileHandler running Get*********")

		data, err := users.StructurizeProfileLists(r)
		if err != nil {
			fmt.Println("StructurizeProfileLists ERR: ", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/profile/image":
		imageName, err := structs.GetUserImage(r)
		if err != nil {
			return
		}
		structs.ServeImage(w, "./media/userImages/", imageName)
	case "/postImage":
		postType := r.URL.Query().Get("postType")
		postID, err := strconv.Atoi(r.URL.Query().Get("postID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		if postType == "group" {
			post, err := structs.GetGroupPost(postID)
			if err != nil {
				fmt.Println("err in getting { POST }")
				return
			}
			structs.ServeImage(w, "./media/postImages/", post.Image.Name)
		} else {
			post, err := structs.GetPost(postID)
			if err != nil {
				fmt.Println("err in getting { POST }")
				return
			}
			structs.ServeImage(w, "./media/postImages/", post.Image.Name)
		}

	case "/commentImage":
		fmt.Println("/comment")
		commentType := r.URL.Query().Get("commentType")
		postID, err := strconv.Atoi(r.URL.Query().Get("postID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		commentID, err := strconv.Atoi(r.URL.Query().Get("commentID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		if commentType == "group" {
			comment, err := structs.GetGroupComment(postID, commentID)
			if err != nil {
				fmt.Println("getComment ERR: ", err)
				http.Error(w, "Invalid user ID.", http.StatusBadRequest)
				return
			}
			structs.ServeImage(w, "./media/commentImages/", comment.Image.Name)
		} else {
			comment, err := structs.GetComment(postID, commentID)
			if err != nil {
				fmt.Println("getComment ERR: ", err)
				http.Error(w, "Invalid user ID.", http.StatusBadRequest)
				return
			}

			structs.ServeImage(w, "./media/commentImages/", comment.Image.Name)
		}

	default:
		fmt.Println("imageHandler DEFAULT")
		http.Error(w, "404 method not allowed", http.StatusNotFound)
	}
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/message" {
		http.Error(w, "404 method not allowed", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// fmt.Println("*******messagesHandler running*********")

	senderID := r.URL.Query().Get("senderID")
	receiverID := r.URL.Query().Get("receiverID")

	messages, err := structs.GetMessages(senderID, receiverID)
	if err != nil {
		http.Error(w, "500 error getting messages", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func ratingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/rating" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// fmt.Println("**********ratingsHandler running**********")

	var postLikeDislikeInfo structs.RatingInfo

	err := json.NewDecoder(r.Body).Decode(&postLikeDislikeInfo)
	if err != nil {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	postLikeDislikeInfo.UserID = structs.GetUserIDfromSessions(r)

	err = structs.AddEventRating(
		postLikeDislikeInfo.EventID,
		postLikeDislikeInfo.UserID,
		postLikeDislikeInfo.RatingType,
	)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	going, notGoing, err := structs.GetEventRatingsCount(postLikeDislikeInfo.EventID)
	if err != nil {
		fmt.Println("88Main: ", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	var ratingsInfo = structs.RatingDetails{
		Going:    going,
		NotGoing: notGoing,
	}

	resp, err := json.Marshal(ratingsInfo)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/feed" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		// fmt.Println("****************************postsPOSTHandler*************************")
		entryType := r.URL.Query().Get("entryType")

		fmt.Println("{ ENTRYTYPE }: ", entryType)
		userID := structs.GetUserIDfromSessions(r)
		entryID, err := structs.CreateFeedEntry(userID, entryType, r.Body)
		if err != nil {
			return
		}

		resp, err := json.Marshal(structs.EntryID{
			EntryID: entryID,
		})
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else if r.Method == "GET" {
		// fmt.Println("****************************postsGetHandler*************************")
		mainData, err := structs.StructurizeMainData(r)
		if err != nil {
			fmt.Println("structurizeMainDataErrMain: ", err)
			return
		}

		resp, err := json.Marshal(mainData)
		if err != nil {
			fmt.Println("marshal err 326: ", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = database.DBase.Exec("DELETE FROM sessionIDs WHERE sessionID = ?", sessionID.Value)
	if err != nil {
		fmt.Println(err)
	}

	var userLoggedOut bool = true

	resp, err := json.Marshal(userLoggedOut)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userLoginInfo structs.UserLoginRequest

	err := json.NewDecoder(r.Body).Decode(&userLoginInfo)
	if err != nil {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	var resp = structs.LoginResponse{
		IsLoggedIn: false,
	}

	if users.Login(database.DBase, userLoginInfo.Email, userLoginInfo.Password) {
		user, err := structs.GetUserFromEmail(userLoginInfo.Email)
		if err != nil {
			log.Println(err)
			return
		}
		sessions.SetSessionCookie(w, user.UserID)
		resp = structs.LoginResponse{
			IsLoggedIn: true,
			UserID:     user.UserID,
			Nickname:   user.Nickname,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Msg:        "Sign in successful",
		}
	} else {
		resp = structs.LoginResponse{
			Msg: "Sign in failed! Check that your Email or Password is correct",
		}
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var userRegisterInfo structs.User
	err := json.NewDecoder(r.Body).Decode(&userRegisterInfo)
	if err != nil {
		fmt.Println("err in { DECODER }: ", err)
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	var emailExists bool
	err = database.DBase.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", userRegisterInfo.Email).Scan(&emailExists)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	if emailExists {
		var resp = structs.RegisterResponse{
			CreateAccount: false,
			Msg:           "Reistration unsuccessful: Email already exists",
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	passwordHash, err := users.GetHashPassword(userRegisterInfo.PasswordHash)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	userRegisterInfo.PasswordHash = passwordHash

	// send the users info to db
	err = structs.AddUser(userRegisterInfo)
	if err != nil {
		fmt.Println("2222")
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	var resp = structs.RegisterResponse{
		CreateAccount: true,
		Msg:           "Registration successful",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func cookieLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cookieLogin" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		isLoggedIn := sessions.IsLoggedIn(r)
		type resp struct {
			Bool         bool `json:"bool"`
			structs.User `json:"resp"`
		}
		if isLoggedIn { //if has active session in db
			userID := structs.GetUserIDfromSessions(r)
			user, err := structs.GetUser(userID)
			if err != nil {
				fmt.Println("Error in GetUser() in cookieLogin handler", err)
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}
			respUser := structs.User{
				UserID:    user.UserID,
				Nickname:  user.Nickname,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}
			dto := resp{
				Bool: true,
				User: respUser,
			}
			frontendResp, err := json.Marshal(dto)
			if err != nil {
				fmt.Println("Error in marshal() in cookieLogin handler", err)
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(frontendResp)
		} else {
			dto := resp{
				Bool: false,
			}
			frontendResp, err := json.Marshal(dto)
			if err != nil {
				fmt.Println("Error in marshal() in cookieLogin handler", err)
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(frontendResp)
		}
	}
}

func groupMessagesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*************************groupMessagesHandler*********************") //////////
	if r.URL.Path != "/groupMessage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		groupID, err := strconv.Atoi(r.URL.Query().Get("groupID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		groupMessages, err := structs.GetGroupMessages(groupID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(groupMessages)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func groupHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("*************************groupHandler*********************") ////////////
	if r.URL.Path != "/group" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		groupID, err := strconv.Atoi(r.URL.Query().Get("groupID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		group, err := structs.GetGroup(groupID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		posts, err := structs.GetAllPostsFromGroup(group.GroupID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var membersList []int
		err = json.Unmarshal(group.Members, &membersList)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		type respObject struct {
			User        []structs.User         `json:"user"`
			GroupPosts  []structs.GroupPost    `json:"groupPosts"`
			Followers   []structs.FollowStruct `json:"followers"`
			Events      []structs.Event        `json:"events"`
			Description string                 `json:"description"`
		}

		memberInfoStruct := make([]structs.User, 0)

		for _, memberID := range membersList {
			memberInfo, err := structs.GetUser(memberID)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
			var resp structs.User
			resp.FirstName = memberInfo.FirstName
			resp.LastName = memberInfo.LastName
			resp.UserID = memberInfo.UserID
			memberInfoStruct = append(memberInfoStruct, resp)
		}

		followersList, _, err := structs.GetFollowLists(structs.GetUserIDfromSessions(r))
		if err != nil {
			fmt.Println("GetFollowersLists err : ", err)
			return
		}

		followersArr := make([]structs.FollowStruct, 0)

		// Populate profile.Followers
		for _, followerID := range followersList {
			user, err := structs.GetUser(followerID)
			if err != nil {
				fmt.Println(err)
				return
			}

			followersArr = append(followersArr, structs.FollowStruct{
				UserID: user.UserID,
				Name:   user.FirstName + " " + user.LastName,
			})
		}

		events, err := structs.GetAllEventsFromGroupID(group.GroupID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var dto respObject
		dto.User = memberInfoStruct
		dto.GroupPosts = posts
		dto.Followers = followersArr
		dto.Description = group.Description
		dto.Events = events

		resp, err := json.Marshal(dto)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else if r.Method == "POST" {
		entryType := r.URL.Query().Get("entryType")
		userID := structs.GetUserIDfromSessions(r)
		entryID, err := structs.CreateGroupEntry(userID, entryType, r.Body)
		if err != nil {
			return
		}
		resp, err := json.Marshal(structs.EntryID{
			EntryID: entryID,
		})
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("*************************groupsHandler*********************") ////////////////////////
	if r.URL.Path != "/groups" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		allGroups, err := structs.GetAllGroups()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		//filter out groups that the user has already joined
		userID := structs.GetUserIDfromSessions(r)
		joinedGroupArr, err := structs.GetJoinedGroupIDArray(userID)
		if err != nil {
			fmt.Println(err)
			return
		}

		notJoinedGroups := make([]structs.Group, 0)
		if len(joinedGroupArr) != 0 { // if there are joined Groups
			for _, group := range allGroups {
				for index, joinedGroupID := range joinedGroupArr {
					if group.GroupID == joinedGroupID {
						break
					} else if index == len(joinedGroupArr)-1 {
						notJoinedGroups = append(notJoinedGroups, group)
					}
				}
			}
		} else {
			notJoinedGroups = allGroups
		}

		resp, err := json.Marshal(notJoinedGroups)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createGroup" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		var newGroupInfo structs.Group
		err := json.NewDecoder(r.Body).Decode(&newGroupInfo)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
		members := []int{newGroupInfo.CreatorID}
		groupID, err := structs.AddGroup(newGroupInfo.CreatorID, newGroupInfo.Name, newGroupInfo.Description, members)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		groupIDArray, err := structs.GetJoinedGroupIDArray(newGroupInfo.CreatorID)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		groupIDArray = append(groupIDArray, groupID)

		err = structs.UpadateJoinedGroupIDArray(groupIDArray, newGroupInfo.CreatorID)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal("groupCreated")
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	} else {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notifications" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		// fmt.Println("*************************notificationsHandlerGET*********************") ////////////////////////
		userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		allNotifications, err := structs.GetAllNotifications(userID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		if len(allNotifications) == 0 {
			allNotifications = make([]structs.NewNotificationEvent, 0)
		}
		resp, err := json.Marshal(allNotifications)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else if r.Method == "POST" {
		// fmt.Println("*************************notificationsHandlerPOST*********************") ////////////////////////

		var notificationInfo structs.Notification

		err := json.NewDecoder(r.Body).Decode(&notificationInfo)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		err = users.InsertIntoSenderDB(
			notificationInfo.SenderID,
			notificationInfo.RecieverID,
		)
		if err != nil {
			fmt.Println("err in Inserting into { SenderDB }: ", err)
			return
		}

		err = users.InsertIntoRecieverDB(
			notificationInfo.SenderID,
			notificationInfo.RecieverID,
		)
		if err != nil {
			fmt.Println("err in Inserting into { RecieverDB }: ", err)
			return
		}

		resp, err := json.Marshal("{ ADDED TO DB }")
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

var CorsMux *CorsHandler

func main() {
	database.ConnDB()
	defer database.CloseDB()

	setUpHandlers()

	fmt.Println("Starting server at: http://localhost:8080")
	listenServeErr := (http.ListenAndServe(":8080", CorsMux))
	if listenServeErr != nil {
		log.Fatal("ListenAndServe: ", listenServeErr)
	}
}

type CorsHandler struct {
	*http.ServeMux
}

func (c CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")                                      // origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")                                        // methods
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Content-Type, Authorization") // headers
	w.Header().Set("Access-Control-Allow-Credentials", "true")                                                  // cookies

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	c.ServeMux.ServeHTTP(w, r)
}

func setUpHandlers() {
	mux := http.NewServeMux()
	CorsMux = &CorsHandler{ServeMux: mux}

	GlobalManager := users.NewManager()

	CorsMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		http.HandlerFunc(GlobalManager.ServeWS).ServeHTTP(w, r)
	})

	CorsMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mux mux mux"))
	})

	//other handlers
	CorsMux.HandleFunc("/feed", postsHandler)
	CorsMux.HandleFunc("/postImage", imageHandler)
	CorsMux.HandleFunc("/commentImage", imageHandler)
	CorsMux.HandleFunc("/register", registerPageHandler)
	CorsMux.HandleFunc("/login", loginPageHandler)
	CorsMux.HandleFunc("/logout", logoutHandler)
	CorsMux.HandleFunc("/rating", ratingsHandler)
	CorsMux.HandleFunc("/message", messagesHandler)
	CorsMux.HandleFunc("/groupMessage", groupMessagesHandler)
	CorsMux.HandleFunc("/groups", groupsHandler)
	CorsMux.HandleFunc("/group", groupHandler)
	CorsMux.HandleFunc("/profile", profileDataHandler)
	CorsMux.HandleFunc("/profile/image", imageHandler)
	CorsMux.HandleFunc("/changeImage", imageHandler)
	CorsMux.HandleFunc("/notifications", notificationsHandler)
	CorsMux.HandleFunc("/createGroup", createGroupHandler)
	CorsMux.HandleFunc("/cookieLogin", cookieLogin)
}
