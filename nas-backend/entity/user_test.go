package entity

import (
	"fmt"
	"nas-backend/utils"
	"testing"
)

func TestAddUser(t *testing.T) {
	user := User{
		Username: "admin",
		Password: "123456",
		Email:    "user@example.com",
	}
	fmt.Println(ValidLoginInfo(utils.DB, user))
}

func TestGetUserIdByUsername(t *testing.T) {
	userId := GetUserIdByUsername(utils.DB, "admin")
	fmt.Println(userId)
}
