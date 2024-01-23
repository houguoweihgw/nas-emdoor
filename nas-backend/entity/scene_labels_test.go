package entity

import (
	"fmt"
	"nas-backend/utils"
	"testing"
)

func TestGetSceneLabelsByUserID(t *testing.T) {
	utils.InitDB()
	labels, err := GetSceneLabelsByUserID(utils.DB, 16)
	if err != nil {
		return
	}
	fmt.Println(labels)
}
