package users

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(db *sql.DB, email string, password string) bool {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ErrNoRows")
			return false
		}
		fmt.Println(err)
	}
	if CheckPasswordHash(password, dbPassword) {
		return true
	}
	return false
}
