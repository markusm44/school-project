package structs

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	database "social-network/backend/database"
	"strconv"
	"strings"
)

// Add chat list order array
func AddChatListOrderArray(arr []int, userID int) error {
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("INSERT INTO chatListOrderArray (`array`, `userID`) VALUES (?, ?)", data, userID)
	if err != nil {
		return err
	}
	return nil
}

func AddFollowersListArrays(arr []int, userID int) error {
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("INSERT INTO followers (`userID`, `followerIDs`, `followedUsers`) VALUES (?, ?, ?)", userID, data, data)
	if err != nil {
		return err
	}
	return nil
}

func GetChatListOrderArray(userID int) ([]byte, error) {
	var arr []byte
	err := database.DBase.QueryRow("SELECT array FROM chatListOrderArray WHERE userID = ?", userID).Scan(
		&arr,
	)
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func UpdateChatListOrderArray(arr []int, userID int) error {
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("UPDATE chatListOrderArray SET array = ? WHERE userID = ?", data, userID)
	if err != nil {
		return err
	}
	return nil
}

// Get all usersID and username from DB users table
func GetAllUsers() ([]User, error) {
	var users []User

	rows, err := database.DBase.Query("SELECT * FROM users")
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserID,
			&user.Nickname,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PasswordHash,
			// &user.Image,
			&user.Image.Name,
			&user.IsProfilePublic,
			&user.DateOfBirth,
			&user.AboutMe,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Gets post from given postID
func GetGroupPost(PostID int) (GroupPost, error) {
	var groupPost GroupPost
	err := database.DBase.QueryRow("SELECT * FROM groupPosts WHERE postID = ?", PostID).Scan(
		&groupPost.PostID,
		&groupPost.UserID,
		&groupPost.GroupID,
		&groupPost.Content,
		&groupPost.Created,
		&groupPost.Image.Name,
		&groupPost.Fullname,
	)
	if err != nil {
		fmt.Print("BIG ERR { GET POST }: ", err)
		return groupPost, err
	}

	groupPost.CommentsInfo, err = GetAllGroupPostComments(groupPost.PostID, groupPost.GroupID)
	if err != nil {
		fmt.Println("error Getting all comments: ", err)
		return groupPost, err
	}

	return groupPost, nil
}

// Gets post from given postID
func GetPost(PostID int) (Post, error) {
	var post Post
	AlmostPrivateUsers := make([]byte, 0)
	err := database.DBase.QueryRow("SELECT * FROM posts WHERE postID = ?", PostID).Scan(
		&post.PostID,
		&post.UserID,
		&post.Content,
		&post.Created,
		&post.Image.Name,
		&post.Fullname,
		&post.Status,
		&AlmostPrivateUsers,
	)
	if err != nil {
		fmt.Print("BIG ERR { GET POST }: ", err)
		return post, err
	}

	if len(AlmostPrivateUsers) == 0 {
		// Handle the case where the byte slice is empty
		post.AlmostPrivateUsers = []int{}
	} else {
		err = json.Unmarshal(AlmostPrivateUsers, &post.AlmostPrivateUsers)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return post, err
		}
	}

	post.CommentsInfo, err = GetAllComments(post.PostID)
	if err != nil {
		fmt.Println("error Getting all comments: ", err)
		return post, err
	}
	return post, nil
}

func GetComment(postID int, commentID int) (Comments, error) {
	var comment Comments
	err := database.DBase.QueryRow("SELECT * FROM comments WHERE postID = ? AND commentID = ?", postID, commentID).Scan(
		&comment.CommentID,
		&comment.PostID,
		&comment.Content,
		&comment.UserID,
		&comment.Fullname,
		&comment.Created,
		&comment.Image.Name,
	)
	if err != nil {
		fmt.Print("BIG ERR { GET POST }: ", err)
		return comment, err
	}

	return comment, nil
}

func GetGroupComment(postID int, commentID int) (GroupComments, error) {
	var groupComment GroupComments
	err := database.DBase.QueryRow("SELECT * FROM groupPostComments WHERE postID = ? AND commentID = ?", postID, commentID).Scan(
		&groupComment.CommentID,
		&groupComment.GroupID,
		&groupComment.PostID,
		&groupComment.Content,
		&groupComment.UserID,
		&groupComment.Fullname,
		&groupComment.Created,
		&groupComment.Image.Name,
	)
	if err != nil {
		fmt.Print("BIG ERR { GET POST }: ", err)
		return groupComment, err
	}

	return groupComment, nil
}

// Gets all posts for given subforumID
func GetAllPosts() ([]Post, error) {
	posts := make([]Post, 0)
	rows, err := database.DBase.Query("SELECT * FROM posts")
	if err != nil {
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post Post
		AlmostPrivateUsers := make([]byte, 0)
		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.Created,
			&post.Image.Name,
			&post.Fullname,
			&post.Status,
			&AlmostPrivateUsers,
		)
		if err != nil {
			fmt.Println("getPostsErrMethods: ", err)
			return posts, err
		}

		if len(AlmostPrivateUsers) == 0 {
			// Handle the case where the byte slice is empty
			post.AlmostPrivateUsers = []int{}
		} else {
			err = json.Unmarshal(AlmostPrivateUsers, &post.AlmostPrivateUsers)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return posts, err
			}
		}

		post.NumberofCommentsOnPost, err = GetNumberOfComments(post.PostID)
		if err != nil {
			return posts, err
		}

		post.CommentsInfo, err = GetAllComments(post.PostID)
		if err != nil {
			fmt.Println("error Getting all comments: ", err)
			return posts, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPostsFromGroup(groupID int) ([]GroupPost, error) {
	posts := make([]GroupPost, 0)
	rows, err := database.DBase.Query("SELECT * FROM groupPosts WHERE groupID = ?", groupID)
	if err != nil {
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post GroupPost
		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.GroupID,
			&post.Content,
			&post.Created,
			&post.Image.Name,
			&post.Fullname,
		)
		if err != nil {
			fmt.Println("getPostsErrMethods: ", err)
			return posts, err
		}

		post.CommentsInfo, err = GetAllGroupPostComments(post.PostID, groupID)
		if err != nil {
			fmt.Println("error Getting all comments: ", err)
			return posts, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllPostsFromUserID(userID int) ([]Post, error) {
	var posts []Post

	rows, err := database.DBase.Query("SELECT * FROM posts WHERE userID = ?", userID)
	if err != nil {
		return posts, err
	}

	defer rows.Close()
	for rows.Next() {
		var post Post
		AlmostPrivateUsers := make([]byte, 0)
		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.Created,
			&post.Image.Name,
			&post.Fullname,
			&post.Status,
			&AlmostPrivateUsers,
		)
		if err != nil {
			fmt.Println("getPostsErrMethods: ", err)
			return posts, err
		}

		if len(AlmostPrivateUsers) == 0 {
			// Handle the case where the byte slice is empty
			post.AlmostPrivateUsers = []int{}
		} else {
			err = json.Unmarshal(AlmostPrivateUsers, &post.AlmostPrivateUsers)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return posts, err
			}
		}

		post.NumberofCommentsOnPost, err = GetNumberOfComments(post.PostID)
		if err != nil {
			return posts, err
		}

		post.CommentsInfo, err = GetAllComments(post.PostID)
		if err != nil {
			fmt.Println("error Getting all comments: ", err)
			return posts, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetAllEvents() ([]Event, error) {
	events := make([]Event, 0)

	rows, err := database.DBase.Query("SELECT * FROM events")
	if err != nil {
		return events, err
	}

	defer rows.Close()
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.GroupID,
			&event.Content,
			&event.Created,
			&event.Title,
			&event.Date,
		)
		if err != nil {
			fmt.Println("getPostsErrMethods: ", err)
			return events, err
		}

		event.Going, event.NotGoing, err = GetEventRatingsCount(event.EventID)
		if err != nil {
			fmt.Println("EventGoing/NotGoingErr: ", err)
			return events, err
		}

		events = append(events, event)
	}
	return events, nil
}

func GetAllEventsFromGroupID(groupID int) ([]Event, error) {
	events := make([]Event, 0)

	rows, err := database.DBase.Query("SELECT * FROM events WHERE groupID= ?", groupID)
	if err != nil {
		return events, err
	}

	defer rows.Close()
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.GroupID,
			&event.Content,
			&event.Created,
			&event.Title,
			&event.Date,
		)
		if err != nil {
			fmt.Println("getPostsErrMethods: ", err)
			return events, err
		}

		event.Going, event.NotGoing, err = GetEventRatingsCount(event.EventID)
		if err != nil {
			fmt.Println("EventGoing/NotGoingErr: ", err)
			return events, err
		}

		events = append(events, event)
	}
	return events, nil
}

// Write user to DB
func AddUser(user User) error {
	image, err := ImageNameHandler(user.Image.Name, user.Image.Base64String, "./media/userImages/")
	if err != nil {
		fmt.Println("error in ImageNameHandler")
		return err
	}

	fmt.Println("user added to DB")

	// Insert user information into the database
	_, err = database.DBase.Exec("INSERT INTO users (`nickname`, `email`, `firstName`, `lastName`, `password`, `image`, `dateOfBirth`, `aboutMe`) VALUES (?,?,?,?,?,?,?,?)", user.Nickname, user.Email, user.FirstName, user.LastName, user.PasswordHash, image, user.DateOfBirth, user.AboutMe)
	if err != nil {
		log.Println("Error in AddUser: ", err)
		return err
	}

	// Initialize different arrays in the database on user registration
	err = addEmptyArrays(user.Email)
	if err != nil {
		log.Println("Error in addEmptyArrays: ", err)
		return err
	}

	return nil
}

// writes empty arrays for different data
func addEmptyArrays(email string) error {
	//get the use struct
	user, err := GetUserFromEmail(email)
	if err != nil {
		log.Println("Error in GetUserFromNicknameOrEmail: ", err)
		return err
	}
	//take userID for insertion into db
	id := user.UserID
	// initialize into db
	err = AddChatListOrderArray([]int{}, id)
	if err != nil {
		log.Println("Error in AddChatListOrderArray: ", err)
		return err
	}
	err = AddJoinedGroupIDArray([]int{}, id)
	if err != nil {
		log.Println("Error in AddJoinedGroupIDArray: ", err)
		return err
	}
	err = AddFollowersListArrays([]int{}, id)
	if err != nil {
		log.Println("Error in AddFollowersListArrays: ", err)
		return err
	}
	return nil
}

// Get user info from username or email
func GetUserFromEmail(email string) (User, error) {
	var user User
	row := database.DBase.QueryRow("SELECT userID, firstName, lastName, email, password FROM users WHERE email = ?", email)
	err := row.Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Get user info from userID
func GetUser(userID int) (User, error) {
	var user User
	row := database.DBase.QueryRow("SELECT * FROM users WHERE userID = ?", userID)
	err := row.Scan(
		&user.UserID,
		&user.Nickname,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.Image.Name,
		&user.IsProfilePublic,
		&user.DateOfBirth,
		&user.AboutMe,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Gets UserID from Sessions table in DB
func GetUserIDfromSessions(r *http.Request) int {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Println("error reading cookie: ", err)
		return -1
	}
	var userID int
	rows := database.DBase.QueryRow("SELECT userID FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&userID)
	if rows != nil {
		fmt.Println(rows)
	}
	return userID
}

// Add group message info into DB
func AddGroupMessages(senderID int, senderName string, groupID int, content string, date string) error {
	_, err := database.DBase.Exec("INSERT INTO groupMessages (`senderID`, `senderName`, `groupID`, `content`,`date`) VALUES (?, ?, ?, ?, ?)", senderID, senderName, groupID, content, date)
	if err != nil {
		return err
	}
	return nil
}

// Get group message from DB
func GetGroupMessages(groupID int) ([]GroupMessages, error) {
	messages := make([]GroupMessages, 0)
	rows, err := database.DBase.Query("SELECT * FROM groupMessages WHERE groupID = ?", groupID)
	if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var message GroupMessages
		err := rows.Scan(
			&message.GroupMessageID,
			&message.SenderID,
			&message.SenderName,
			&message.GroupID,
			&message.Content,
			&message.Date,
		)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// Get Joined GroupID Array from userID
func GetJoinedGroupIDArray(userID int) ([]int, error) {
	var arr []byte
	res := make([]int, 0)
	err := database.DBase.QueryRow("SELECT array FROM joinedGroupIDArray WHERE userID =?", userID).Scan(
		&arr,
	)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(arr, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Add Joined GroupID Array
func AddJoinedGroupIDArray(arr []int, userID int) error {
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("INSERT INTO joinedGroupIDArray (`array`, `userID`) VALUES (?, ?)", data, userID)
	if err != nil {
		return err
	}
	return nil
}

func UpadateJoinedGroupIDArray(arr []int, userID int) error {
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("UPDATE joinedGroupIDArray SET array = ? WHERE userID = ?", data, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllGroups() ([]Group, error) {
	groups := make([]Group, 0)
	rows, err := database.DBase.Query("SELECT * FROM groups")
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group Group
		var members, created string
		err := rows.Scan(
			&group.GroupID,
			&group.CreatorID,
			&group.Name,
			&group.Description,
			&members,
			&created,
		)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func GetGroup(groupID int) (Group, error) {
	var group Group
	err := database.DBase.QueryRow("SELECT * FROM groups WHERE groupID =?", groupID).Scan(
		&group.GroupID,
		&group.CreatorID,
		&group.Name,
		&group.Description,
		&group.Members,
		&group.Created,
	)
	if err != nil {
		return group, err
	}
	return group, nil
}

func AddGroup(creatorID int, name string, description string, membersArr []int) (int, error) {
	data, err := json.Marshal(membersArr)
	if err != nil {
		return 0, err
	}
	res, err := database.DBase.Exec("INSERT INTO groups (`creatorID`, `name`, `description`, `members`) VALUES (?,?, ?, ?)", creatorID, name, description, data)
	if err != nil {
		return 0, err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(groupID), nil
}

func UpdateGroupMembers(membersArr []int, groupID int) error {
	data, err := json.Marshal(membersArr)
	if err != nil {
		return err
	}
	_, err = database.DBase.Exec("UPDATE groups SET members = ? WHERE groupID = ?", data, groupID)
	if err != nil {
		return err
	}
	return nil
}

// Add message info into DB
func AddMessages(senderID int, receiverID int, content string, date string) error {
	_, err := database.DBase.Exec("INSERT INTO messages (`senderID`, `receiverID`, `content`,`date`) VALUES (?, ?, ?, ?)", senderID, receiverID, content, date)
	if err != nil {
		return err
	}
	return nil
}

// Get message from DB
func GetMessages(senderID string, receiverID string) ([]Messages, error) {
	var messages []Messages
	rows, err := database.DBase.Query("SELECT * FROM messages WHERE senderID = ? AND receiverID = ? OR senderID = ? AND receiverID = ?", senderID, receiverID, receiverID, senderID)
	if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var message Messages
		err := rows.Scan(
			&message.MessageID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.Date,
		)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func GetAllComments(postID int) ([]Comments, error) {
	comments := make([]Comments, 0)
	rows, err := database.DBase.Query("SELECT * FROM comments WHERE postID = ?", postID)
	if err != nil {
		return comments, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comments
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.Content,
			&comment.UserID,
			&comment.Fullname,
			&comment.Created,
			&comment.Image.Name,
		)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}

func AddGroupPost(userID int, postDetails GroupPost) error {
	image, err := ImageNameHandler(postDetails.Image.Name, postDetails.Image.Base64String, "./media/postImages/")
	if err != nil {
		fmt.Println("error in ImageNameHandler")
		return err
	}
	_, err = database.DBase.Exec("INSERT INTO groupPosts (`userID`, `groupID`, `content`, `created`, `image`, `fullName`) VALUES (?, ?, ?, ?, ?, ?)", userID, postDetails.GroupID, postDetails.Content, postDetails.Created, image, postDetails.Fullname)
	if err != nil {
		fmt.Println("inAddPostErr: ", err)
		return err
	}
	return nil
}

func AddGroupComment(userID int, groupcommentDetails GroupComments) error {
	image, err := ImageNameHandler(groupcommentDetails.Image.Name, groupcommentDetails.Image.Base64String, "./media/commentImages/")
	if err != nil {
		fmt.Println("error in ImageNameHandler")
		return err
	}

	_, err = database.DBase.Exec("INSERT INTO groupPostComments (`groupID`, `postID`, `content`, `userID`, `fullName`, `created`, `image`) VALUES (?,?,?,?,?,?,?)", groupcommentDetails.GroupID, groupcommentDetails.PostID, groupcommentDetails.Content, userID, groupcommentDetails.Fullname, groupcommentDetails.Created, image)
	if err != nil {
		return err
	}
	return nil
}

func GetAllGroupPostComments(postID, groupID int) ([]GroupComments, error) {
	comments := make([]GroupComments, 0)
	rows, err := database.DBase.Query("SELECT * FROM groupPostComments WHERE postID = ? AND groupID = ?", postID, groupID)
	if err != nil {
		return comments, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment GroupComments
		err := rows.Scan(
			&comment.CommentID,
			&comment.GroupID,
			&comment.PostID,
			&comment.Content,
			&comment.UserID,
			&comment.Fullname,
			&comment.Created,
			&comment.Image.Name,
		)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}

func AddEventRating(eventID int, userID int, ratingType string) error {
	// Check if the user has already rated the event
	row := database.DBase.QueryRow("SELECT rating FROM ratings WHERE eventID = ? AND userID = ?", eventID, userID)

	var existingRating string
	err := row.Scan(&existingRating)
	if err != nil {
		// User has not rated this event before, insert a new rating
		_, err := database.DBase.Exec("INSERT INTO ratings (`eventID`, `userID`, `rating`) VALUES (?,?,?)", eventID, userID, ratingType)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	if existingRating == ratingType {
		// User wants to remove their rating
		DeleteRating(userID, eventID)
	} else {
		// User wants to update their rating
		UpdateRating(userID, eventID, ratingType)
	}

	return nil
}

// return likes, dislikes, error
func GetEventRatingsCount(eventID int) (int, int, error) {
	//counts going
	var going int
	err := database.DBase.QueryRow("SELECT COUNT(*) FROM ratings WHERE eventID = ? AND rating = ?", eventID, "going").Scan(&going)
	if err != nil {
		return 0, 0, err
	}

	//counts notGoing
	var notGoing int
	err = database.DBase.QueryRow("SELECT COUNT(*) FROM ratings WHERE eventID = ? AND rating = ?", eventID, "notGoing").Scan(&notGoing)
	if err != nil {
		return 0, 0, err
	}

	return going, notGoing, nil
}

func UpdateRating(userID int, eventID int, ratingType string) error {
	_, err := database.DBase.Exec("UPDATE ratings SET rating = ? WHERE eventID = ? AND userID = ?", ratingType, eventID, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRating(userID int, eventID int) error {
	_, err := database.DBase.Exec("DELETE FROM ratings WHERE eventID = ? AND userID = ?", eventID, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetNumberOfComments(postID int) (int, error) {
	var count int
	err := database.DBase.QueryRow("SELECT COUNT(*) FROM comments WHERE postID = ?", postID).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func GetCommentRating(commentID int) (int, int, error) {
	//counts likes
	var likes int
	err := database.DBase.QueryRow("SELECT COUNT(*) FROM ratings WHERE commentID = ? AND rating = ?", commentID, "like").Scan(&likes)
	if err != nil {
		return likes, 0, err
	}
	//counts dislikes
	var dislikes int
	err = database.DBase.QueryRow("SELECT COUNT(*) FROM ratings WHERE commentID = ? AND rating = ?", commentID, "dislike").Scan(&dislikes)
	if err != nil {
		return likes, dislikes, err
	}
	return likes, dislikes, nil
}

func GetProfilePosts(recieverID int, idFromProfilePageURL int) ([]Post, error) {
	AccessiblePosts := make([]Post, 0)
	// Call GetAllPosts to get the list of posts
	posts, err := GetAllPostsFromUserID(idFromProfilePageURL)
	if err != nil {
		return make([]Post, 0), err
	}
	fmt.Print("\n\n{ ALL POSTS } ShowPostBasedOnStatus: ", posts, "\n\n")
	// Iterate through the retrieved posts
	for _, post := range posts {
		fmt.Println("{ POST DATA TYPE }: ", reflect.TypeOf(post.AlmostPrivateUsers))
		// Check the post's status and handle it accordingly
		switch post.Status {
		case "public":
			// Handle public posts
			AccessiblePosts = append(AccessiblePosts, post)
		case "private":
			// Handle private posts
			if post.UserID == recieverID {
				// Check if the current user made the post
				AccessiblePosts = append(AccessiblePosts, post)
			}
			followerList, _, err := GetFollowLists(post.UserID)
			if err != nil {
				fmt.Println("GetFollowLists err in RemoveFromFollowedLists: ", err)
				return make([]Post, 0), err
			}
			// Check if recieverID is among the followers
			for _, followerID := range followerList {
				if followerID == recieverID {
					AccessiblePosts = append(AccessiblePosts, post)
					break
				}
			}
		case "almostPrivate":
			fmt.Println("{ ALMOST PRIVATE }")
			// Handle almostPrivate posts
			if post.UserID == recieverID {
				// Check if the current user made the post
				AccessiblePosts = append(AccessiblePosts, post)
			}
			fmt.Println("almostePrivateList { BYTE }: ", post.AlmostPrivateUsers)
			// Check if recieverID is among the almostPrivateUsers
			for _, followerID := range post.AlmostPrivateUsers {
				if followerID == recieverID {
					AccessiblePosts = append(AccessiblePosts, post)
					break
				}
			}
		default:
			fmt.Println("{ DEFAULT } loopERR")
			return make([]Post, 0), fmt.Errorf("Invalid entry type")
		}
	}
	return AccessiblePosts, nil
}

func GetPostsBasedOnPostStatus(recieverID int) ([]Post, error) {
	AccessiblePosts := make([]Post, 0)
	// Call GetAllPosts to get the list of posts
	posts, err := GetAllPosts()
	if err != nil {
		return make([]Post, 0), err
	}
	// Iterate through the retrieved posts
	for _, post := range posts {
		// Check the post's status and handle it accordingly
		switch post.Status {
		case "public":
			// Handle public posts
			AccessiblePosts = append(AccessiblePosts, post)
		case "private":
			// Handle private posts
			if post.UserID == recieverID {
				// Check if the current user made the post
				AccessiblePosts = append(AccessiblePosts, post)
			}
			followerList, _, err := GetFollowLists(post.UserID)
			if err != nil {
				fmt.Println("GetFollowLists err in RemoveFromFollowedLists: ", err)
				return make([]Post, 0), err
			}
			// Check if recieverID is among the followers
			for _, followerID := range followerList {
				if followerID == recieverID {
					AccessiblePosts = append(AccessiblePosts, post)
					break
				}
			}
		case "almostPrivate":
			// Handle almostPrivate posts
			if post.UserID == recieverID {
				// Check if the current user made the post
				AccessiblePosts = append(AccessiblePosts, post)
			}
			// Check if recieverID is among the almostPrivateUsers
			for _, followerID := range post.AlmostPrivateUsers {
				if followerID == recieverID {
					AccessiblePosts = append(AccessiblePosts, post)
					break
				}
			}
		default:
			fmt.Println("{ DEFAULT } loopERR")
			return make([]Post, 0), fmt.Errorf("Invalid entry type")
		}
	}
	return AccessiblePosts, nil
}

func GetGroupPostsForProfile(recieverID int, groupID int) ([]GroupPost, error) {
	AccessiblePosts := make([]GroupPost, 0)
	// Call GetAllPosts to get the list of posts
	groupPosts, err := GetAllPostsFromGroup(groupID)
	if err != nil {
		return make([]GroupPost, 0), err
	}

	for _, groupPost := range groupPosts {
		if groupPost.UserID == recieverID {
			// Check if the current user made the post
			AccessiblePosts = append(AccessiblePosts, groupPost)
		}
		followerList, _, err := GetFollowLists(groupPost.UserID)
		if err != nil {
			fmt.Println("GetFollowLists err in RemoveFromFollowedLists: ", err)
			return make([]GroupPost, 0), err
		}
		// Check if recieverID is among the followers
		for _, followerID := range followerList {
			if followerID == recieverID {
				AccessiblePosts = append(AccessiblePosts, groupPost)
				break
			}
		}
	}

	return AccessiblePosts, nil
}

func GetFollowLists(userID int) ([]int, []int, error) {
	followers := make([]byte, 0)
	followedUsers := make([]byte, 0)

	err := database.DBase.QueryRow("SELECT followerIDs, followedUsers FROM followers WHERE userID = ?", userID).Scan(&followers, &followedUsers)
	if err != nil {
		fmt.Println("ERR in GetFollowersList queryRow: ", err)
		return make([]int, 0), make([]int, 0), err
	}

	followersList := make([]int, 0)
	followedUsersList := make([]int, 0)

	err = json.Unmarshal(followers, &followersList)
	if err != nil {
		fmt.Println("ERR in GetFollowersList 1st Unmarshal: ", err)
		return make([]int, 0), make([]int, 0), err
	}

	err = json.Unmarshal(followedUsers, &followedUsersList)
	if err != nil {
		fmt.Println("ERR in GetFollowersList 2nd Unmarshal: ", err)
		return make([]int, 0), make([]int, 0), err
	}

	return followersList, followedUsersList, nil
}

func RemoveFromFollowLists(userID int, followedUserID int) error {
	_, followedList, err := GetFollowLists(userID)
	if err != nil {
		fmt.Println("GetFollowLists err in RemoveFromFollowedLists: ", err)
		return err
	}

	followerList, _, err := GetFollowLists(followedUserID)
	if err != nil {
		fmt.Println("GetFollowLists err in RemoveFromFollowerLists: ", err)
		return err
	}

	// Find the indices of the users to be removed.
	followedIndex := -1
	followerIndex := -1

	for i, userID := range followedList {
		if userID == followedUserID {
			followedIndex = i
			break
		}
	}

	for i, id := range followerList {
		if id == userID {
			followerIndex = i
			break
		}
	}

	if followedIndex < 0 || followerIndex < 0 {
		// User not found in one of the lists.
		fmt.Println("User not found in one of the lists.")
		return errors.New("User not found in one of the lists")
	}

	// Update the lists.
	followedList = append(followedList[:followedIndex], followedList[followedIndex+1:]...)
	followerList = append(followerList[:followerIndex], followerList[followerIndex+1:]...)

	// Convert lists to JSON.
	followedUsersJSON, err := json.Marshal(followedList)
	if err != nil {
		return fmt.Errorf("JSON Marshal error: %d", err)
	}
	followerUsersJSON, err := json.Marshal(followerList)
	if err != nil {
		return fmt.Errorf("JSON Marshal error: %d", err)
	}

	// Update the database with the modified lists.
	_, err = database.DBase.Exec("UPDATE followers SET followedUsers = ? WHERE userID = ?", followedUsersJSON, userID)
	if err != nil {
		return fmt.Errorf("Database update error: %d", err)
	}

	_, err = database.DBase.Exec("UPDATE followers SET followerIDs = ? WHERE userID = ?", followerUsersJSON, followedUserID)
	if err != nil {
		return fmt.Errorf("Database update error: %d", err)
	}

	return nil
}

func StructurizeMainData(r *http.Request) (MainStruct, error) {
	var mainStruct MainStruct

	userID := GetUserIDfromSessions(r)

	allPosts, err := GetPostsBasedOnPostStatus(userID)
	if err != nil {
		fmt.Println(err)
		return mainStruct, err
	}

	allEvents, err := GetAllEvents()
	if err != nil {
		fmt.Println(err)
		return mainStruct, err
	}

	// Get and format all Groups a user has joined
	userJoinedGroups, err := GetJoinedGroupIDArray(userID)
	if err != nil {
		fmt.Println("error getting GetJoinedGroupIDArray: ", err)
		return mainStruct, err
	}

	userJoinedGroupsResp := make([]Group, 0)
	for _, groupID := range userJoinedGroups {
		joinedGroup, err := GetGroup(groupID)
		if err != nil {
			fmt.Println("error getting GetGroup: ", err)
			return mainStruct, err
		}
		var groupRespObj Group
		groupRespObj.GroupID = joinedGroup.GroupID
		groupRespObj.Name = joinedGroup.Name

		userJoinedGroupsResp = append(userJoinedGroupsResp, groupRespObj)
	}

	fmt.Println(userJoinedGroupsResp)

	data := MainStruct{
		AllPosts:         allPosts,
		AllEvents:        allEvents,
		UserJoinedGroups: userJoinedGroupsResp,
	}

	return data, err
}

func CreateFeedEntry(userID int, entryType string, body io.ReadCloser) (int, error) {
	switch entryType {
	case "post":
		var data Post
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		// Create a post entry
		err = AddPost(userID, data)

		// get the id of the post that was last created
		var postID int
		err = database.DBase.QueryRow("SELECT postID FROM posts ORDER BY postID DESC").Scan(&postID)
		if err != nil {
			return 0, err
		}
		return postID, err
	case "comment":
		var data Comments
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		// Create a comment entry
		err = AddComment(userID, data)

		// get the id of the comment that was last created
		var commentID int
		err = database.DBase.QueryRow("SELECT commentID FROM comments ORDER BY commentID DESC").Scan(&commentID)
		if err != nil {
			fmt.Println("error in getting feed commentID: ", err)
			return 0, err
		}
		return commentID, err
	default:
		return 0, fmt.Errorf("invalid entry type")
	}
}
func CreateGroupEntry(userID int, entryType string, body io.ReadCloser) (int, error) {
	switch entryType {
	case "event":
		var data Event
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		// Create an event entry
		err = AddEvent(userID, data)

		// get the id of the event that was last created
		var eventID int
		err = database.DBase.QueryRow("SELECT eventID FROM events ORDER BY eventID DESC").Scan(&eventID)
		if err != nil {
			fmt.Println("dbaseErr { EVENT }: ", err)
			return 0, err
		}
		return eventID, err
	case "post":
		var data GroupPost
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		// Create a post entry
		err = AddGroupPost(userID, data)

		// get the id of the post that was last created
		var postID int
		err = database.DBase.QueryRow("SELECT postID FROM groupPosts ORDER BY postID DESC").Scan(&postID)
		if err != nil {
			return 0, err
		}
		return postID, err
	case "comment":
		var data GroupComments
		err := json.NewDecoder(body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		// Create a comment entry
		err = AddGroupComment(userID, data)

		// get the id of the event that was last created
		var commentID int
		err = database.DBase.QueryRow("SELECT commentID FROM groupPostComments ORDER BY commentID DESC").Scan(&commentID)
		if err != nil {
			fmt.Println("error in getting Group CommentID: ", err)
			return 0, err
		}
		return commentID, err
	default:
		return 0, fmt.Errorf("invalid entry type")
	}
}

// writes subforum to DB
func AddPost(userID int, postDetails Post) error {
	image, err := ImageNameHandler(postDetails.Image.Name, postDetails.Image.Base64String, "./media/postImages/")
	if err != nil {
		fmt.Println("error in ImageNameHandler")
		return err
	}

	JSON, err := json.Marshal(postDetails.AlmostPrivateUsers)
	if err != nil {
		fmt.Println("{ JSON ERR } ")
		return err
	}

	if postDetails.Status == "almostPrivate" {
		fmt.Println("in almostPrivate")
		_, err := database.DBase.Exec("INSERT INTO posts (`userID`, `content`, `created`, `image`, `fullName`, `status`, `almostPrivateUsers`) VALUES (?, ?, ?, ?, ?, ?, ?)", userID, postDetails.Content, postDetails.Created, image, postDetails.Fullname, postDetails.Status, JSON)
		if err != nil {
			fmt.Println("inAddPostErr: ", err)
			return err
		}
	} else {
		fmt.Println("in else")
		_, err := database.DBase.Exec("INSERT INTO posts (`userID`, `content`, `created`, `image`, `fullName`, `status`, `almostPrivateUsers`) VALUES (?, ?, ?, ?, ?, ?, ?)", userID, postDetails.Content, postDetails.Created, image, postDetails.Fullname, postDetails.Status, make([]byte, 0))
		if err != nil {
			fmt.Println("inAddPostErr: ", err)
			return err
		}
	}

	return nil
}

func AddEvent(userID int, eventDetails Event) error {
	_, err := database.DBase.Exec("INSERT INTO events (`userID`, `groupID`, `content`, `created`, `Title`, `date`) VALUES (?, ?, ?, ?, ?, ?)", userID, eventDetails.GroupID, eventDetails.Content, eventDetails.Created, eventDetails.Title, eventDetails.Date)
	if err != nil {
		fmt.Println("inAddEventErr: ", err)
		return err
	}
	return nil
}

func AddComment(userID int, commentDetails Comments) error {
	image, err := ImageNameHandler(commentDetails.Image.Name, commentDetails.Image.Base64String, "./media/commentImages/")
	if err != nil {
		fmt.Println("error in ImageNameHandler")
		return err
	}

	_, err = database.DBase.Exec("INSERT INTO comments (`postID`, `content`, `userID`, `fullName`, `created`, `image`) VALUES (?,?,?,?,?,?)", commentDetails.PostID, commentDetails.Content, userID, commentDetails.Fullname, commentDetails.Created, image)
	if err != nil {
		return err
	}
	return nil
}

// returns added notificationID and error
func AddNotification(subtype, content string, senderID, recieverID int) (int, error) {
	res, err := database.DBase.Exec("INSERT INTO notifications (`notificationSubtype`, `senderID`, `recieverID`, `content`) VALUES (?,?,?,?)", subtype, senderID, recieverID, content)
	if err != nil {
		return 0, err
	}
	notificationID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	notificationIDint := int(notificationID)
	if err != nil {
		return notificationIDint, err
	}
	return notificationIDint, nil
}

func DeleteNotification(notificationID int) error {
	_, err := database.DBase.Exec("DELETE FROM notifications WHERE notificationID = ?", notificationID)
	if err != nil {
		return err
	}
	return nil
}
func GetAllNotifications(recieverID int) ([]NewNotificationEvent, error) {
	var notifications []NewNotificationEvent
	rows, err := database.DBase.Query("SELECT * FROM notifications WHERE recieverID = ?", recieverID)
	if err != nil {
		return notifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification Notification
		err := rows.Scan(
			&notification.NotificationID,
			&notification.NotificationSubtype,
			&notification.SenderID,
			&notification.RecieverID,
			&notification.Content,
			&notification.Created,
		)
		notificationEvent := NewNotificationEvent{
			Subtype: notification.NotificationSubtype,
			Body:    notification,
		}
		if err != nil {
			return notifications, err
		}
		notifications = append(notifications, notificationEvent)
	}
	return notifications, nil
}

func ImageNameHandler(imageName string, Base64String string, savePath string) (string, error) {
	if imageName == "" {
		imageName = "default_image.png"
	} else {
		err := SaveBase64Image(
			Base64String,
			imageName,
			savePath,
		)
		if err != nil {
			fmt.Println("{ SAVE BASE IMAGE } err: ", err)
			return "", err
		}
	}
	return imageName, nil
}

func GetUserImage(r *http.Request) (string, error) {
	var imageName string
	userID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		fmt.Println("Atoi err: ", err)
		return "", err
	}
	err = database.DBase.QueryRow("SELECT image FROM users WHERE userID = ?", userID).Scan(
		&imageName,
	)
	if err != nil {
		fmt.Println("error getting image from user: ", err)
		return "", err
	}
	return imageName, nil
}

func GetPostImage(postID int) (string, error) {
	var imageName string
	err := database.DBase.QueryRow("SELECT image FROM posts WHERE postID = ?", postID).Scan(
		&imageName,
	)
	if err != nil {
		fmt.Println("error getting image from user: ", err)
		return "", err
	}
	return imageName, nil
}

func SaveBase64Image(base64String, imgName, path string) error {
	// Supported image types and their corresponding prefixes
	imageTypes := map[string]string{
		"jpg": "data:image/jpeg;base64,",
		"png": "data:image/png;base64,",
		"gif": "data:image/gif;base64,",
	}

	// Initialize imageType to an empty string
	var imageType string

	// Determine the image type based on the content of the base64 string
	for ext, prefix := range imageTypes {
		if strings.HasPrefix(base64String, prefix) {
			imageType = ext
			base64String = strings.TrimPrefix(base64String, prefix)
			break
		}
	}

	// Check if a valid image type was found
	if imageType == "" {
		return fmt.Errorf("unsupported image type or too large")
	}

	// Decode the base64 string into binary image data
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println("{ MKDIR } error: ", err)
			return err
		}
	}

	// Write the decoded image data to the specified file path
	err = ioutil.WriteFile(path+imgName, decoded, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ServeImage(w http.ResponseWriter, imagePath, imageName string) {
	file, err := os.Open(filepath.Join(imagePath, imageName))
	if err != nil {
		fmt.Println("ERR os.open image: ", err)
		http.Error(w, "Error serving image", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Determine the content type based on the file extension.
	fileExtension := filepath.Ext(imageName)
	switch fileExtension {
	case ".png", ".PNG":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".JPG", ".jpeg", ".JPEG":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif", ".GIF":
		w.Header().Set("Content-Type", "image/gif")
	default:
		fmt.Println("ERR serving image")
		http.Error(w, "Error serving image - Unsupported image format", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println("ERR copy image")
		http.Error(w, "Error serving image", http.StatusInternalServerError)
	}
}
