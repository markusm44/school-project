package sessions

import (
	"fmt"
	"net/http"
	database "social-network/backend/database"
	"time"

	"github.com/gofrs/uuid"
)

func SetSessionCookie(w http.ResponseWriter, userID int) {
	sessionID, _ := uuid.NewV4()
	expiration := time.Now().Add(time.Hour)
	cookie := &http.Cookie{
		Name:    "sessionID",
		Value:   sessionID.String(),
		Expires: expiration,
		Path:    "/",
	}

	http.SetCookie(w, cookie)
	// Delete old cookie
	_, err := database.DBase.Exec("DELETE FROM sessionIDs WHERE userID = ?", userID)
	if err != nil {
		fmt.Println(err)
	}
	_, err = database.DBase.Exec("INSERT INTO sessionIDs(sessionID, expires, userID) VALUES (?,?,?)", sessionID.String(), expiration, userID)
	if err != nil {
		fmt.Println(err)
	}
}

// Get userID from username
func getUserID(username string) string {
	var userID string
	rows := database.DBase.QueryRow("SELECT userID FROM users WHERE username = ? OR email = ?", username, username).Scan(&userID)
	if rows != nil {
		fmt.Println(rows)
	}
	return userID
}

func IsLoggedIn(r *http.Request) bool {
	// if there is no such cookie, return false
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return false
	}
	// SQL will the return the number of matches it got from cookie table
	var count int
	err = database.DBase.QueryRow("SELECT COUNT(*) FROM sessionIDs").Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	//incase the table is empty return false
	if count == 0 {
		return false
	}
	err = database.DBase.QueryRow("SELECT COUNT(*) FROM sessionIDs WHERE sessionID = ?", cookie.Value).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count == 1 {
		return true
	}
	return false
}
