package entity

import (
	"fmt"
	"nas-backend/utils"
	"testing"
)

//func TestInsertAlbumByUserID(t *testing.T) {
//	utils.InitDB()
//	fmt.Println(InsertAlbumByUserID(utils.DB, 1, "深圳打工日记", "2023年7月开始的深圳打工日记"))
//}

func TestQueryAllAlbumsByUserId(t *testing.T) {
	utils.InitDB()
	fmt.Println(QueryAllAlbumsByUserId(utils.DB, 1))
}

func TestDeleteAlbumByAlbumID(t *testing.T) {
	utils.InitDB()
	fmt.Println(DeleteAlbumByAlbumID(utils.DB, 5))
}
