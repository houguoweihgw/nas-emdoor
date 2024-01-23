package entity

import (
	"fmt"
	"nas-backend/utils"
	"testing"
)

func TestInsertPhoto(t *testing.T) {
	utils.InitDB()
	//InsertPhoto(utils.DB)
}

func TestGetAllPhotosByUserID(t *testing.T) {
	utils.InitDB()
	GetAllPhotosByUserID(utils.DB, 1)
}

func TestUpdatePhotoCollectedByPhotoID(t *testing.T) {
	utils.InitDB()
	err := UpdatePhotoCollectedByPhotoID(utils.DB, 4)
	fmt.Println(err)
}

func TestGetAllCollectedPhotosByUserID(t *testing.T) {
	utils.InitDB()
	photos, err := GetAllCollectedPhotosByUserID(utils.DB, 1)
	if err != nil {
		return
	}
	fmt.Println(photos)
}

func TestAddPhotoToAlbum(t *testing.T) {
	utils.InitDB()
	err := AddPhotoToAlbum(utils.DB, 13, 2)
	fmt.Println(err)
}
