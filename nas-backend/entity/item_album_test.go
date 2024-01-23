package entity

import (
	"fmt"
	"nas-backend/utils"
	"testing"
)

func TestGetItemsInAlbum(t *testing.T) {
	utils.InitDB()
	res, err := GetItemsInAlbum(utils.DB, 2)
	fmt.Println(res, err)
}
