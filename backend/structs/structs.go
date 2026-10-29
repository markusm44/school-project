package structs

import (
	"encoding/json"
)

type Subforum struct {
	SubforumID  int    `json:"subforumID"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Post struct {
	PostID                 int        `json:"postID"`
	UserID                 int        `json:"userID"`
	Content                string     `json:"content"`
	Image                  Image      `json:"image"`
	NumberofCommentsOnPost int        `json:"numberofCommentsOnPost"`
	Created                string     `json:"created"`
	Fullname               string     `json:"fullname"`
	Status                 string     `json:"status"`
	AlmostPrivateUsers     []int      `json:"almostPrivateUsers"`
	CommentsInfo           []Comments `json:"commentsInfo"`
}
type GroupPost struct {
	GroupID                int             `json:"groupID"`
	PostID                 int             `json:"postID"`
	UserID                 int             `json:"userID"`
	Content                string          `json:"content"`
	Image                  Image           `json:"image"`
	NumberofCommentsOnPost int             `json:"numberofCommentsOnPost"`
	Created                string          `json:"created"`
	Fullname               string          `json:"fullname"`
	Status                 string          `json:"status"`
	AlmostPrivateUsers     []int           `json:"almostPrivateUsers"`
	CommentsInfo           []GroupComments `json:"commentsInfo"`
}

type Event struct {
	EventID  int    `json:"eventID"`
	UserID   int    `json:"userID"`
	GroupID  int    `json:"groupID"`
	Content  string `json:"content"`
	Title    string `json:"title"`
	Created  string `json:"created"`
	Fullname string `json:"fullname"`
	Going    int    `json:"going"`
	NotGoing int    `json:"notgoing"`
	Date     string `json:"date"`
}

type Comments struct {
	CommentID int    `json:"commentID"`
	PostID    int    `json:"postID"`
	Content   string `json:"content"`
	UserID    int    `json:"userID"`
	Fullname  string `json:"fullname"`
	Created   string `json:"created"`
	Image     Image  `json:"image"`
}

type GroupComments struct {
	CommentID int    `json:"commentID"`
	GroupID   int    `json:"groupID"`
	PostID    int    `json:"postID"`
	Content   string `json:"content"`
	UserID    int    `json:"userID"`
	Fullname  string `json:"fullname"`
	Created   string `json:"created"`
	Image     Image  `json:"image"`
}

type User struct {
	UserID          int    `json:"userID"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	PasswordHash    string `json:"passwordHash"`
	Image           Image  `json:"image"`
	IsProfilePublic bool   `json:"isProfilePublic"`
	DateOfBirth     string `json:"dateOfBirth"`
	AboutMe         string `json:"aboutMe"`
}

type Image struct {
	Base64String string `json:"base64String"`
	Name         string `json:"name"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RatingDetails struct {
	Going    int `json:"going"`
	NotGoing int `json:"notGoing"`
}

type RatingInfo struct {
	RatingType string `json:"ratingType"`
	EventID    int    `json:"eventID"`
	UserID     int    `json:"userID"`
}

type Messages struct {
	MessageID  int    `json:"messageID"`
	SenderID   int    `json:"senderID"`
	ReceiverID int    `json:"receiverID"`
	Content    string `json:"content"`
	Date       string `json:"date"`
}

type RegisterResponse struct {
	CreateAccount bool   `json:"createAccount"`
	Msg           string `json:"msg"`
	ImageState    bool   `json:"imageState"`
}

type LoginResponse struct {
	IsLoggedIn bool   `json:"isLoggedIn"`
	UserID     int    `json:"userID"`
	Nickname   string `json:"nickname"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Msg        string `json:"msg"`
}

type GroupMessages struct {
	GroupMessageID int    `json:"messageID"`
	SenderID       int    `json:"senderID"`
	SenderName     string `json:"senderName"`
	GroupID        int    `json:"groupID"`
	Content        string `json:"content"`
	Date           string `json:"date"`
}

type Group struct {
	GroupID     int             `json:"groupID"`
	CreatorID   int             `json:"creatorID"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Members     json.RawMessage `json:"members"`
	Created     string          `json:"created"`
}

type FollowStatus struct {
	IsFollowing bool `json:"isFollowing"`
	IsRequested bool `json:"isRequested"`
}

type ProfileInfo struct {
	PublicFollowBool bool `json:"publicFollowBool"`
}

type EntryID struct {
	EntryID int `json:"entryID"`
}

type FollowStruct struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
}

type MainStruct struct {
	AllPosts         []Post  `json:"allPosts"`
	AllEvents        []Event `json:"allEvents"`
	UserJoinedGroups []Group `json:"userJoinedGroups"`
}

type GroupStruct struct {
	EntryType string          `json:"entryType"`
	Data      json.RawMessage `json:"entryData"`
}

type ProfileData struct {
	UserStruct User `json:"userStruct"`
	Followers  []FollowStruct
	Followed   []FollowStruct
	Posts      []Post
	GroupPosts []GroupPost
}

type Notification struct {
	NotificationID      int    `json:"notificationID"`
	NotificationSubtype string `json:"subtype"`
	SenderID            int    `json:"senderID"`
	RecieverID          int    `json:"recieverID"`
	Content             string `json:"content"`
	Created             string `json:"created"`
}

type NewNotificationEvent struct {
	Subtype string      `json:"subtype"`
	Body    interface{} `json:"body"`
}
