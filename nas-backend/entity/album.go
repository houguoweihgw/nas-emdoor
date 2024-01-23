package entity

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Album struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	DateCreated time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"dateCreated"`
	UserID      uint      `json:"userId"`
}

// CreateAlbum 根据相册信息在albums表中新建一个相册
func CreateAlbum(db *gorm.DB, album Album) error {
	result := db.Create(&album)
	err := result.Error
	if err != nil {
		log.Printf("albums表插入失败")
		return err
	}
	log.Printf("albums表插入成功, 用户ID：%d, 相册ID：%d, 相册Name：%s", album.UserID, album.ID, album.Name)
	return err
}

// QueryAllAlbumsByUserId 根据用户ID在albums表中查询该用户的所有相册列表
func QueryAllAlbumsByUserId(db *gorm.DB, userID int) ([]Album, error) {
	var albums []Album
	//查询该用户所有相册
	result := db.Where("user_id = ?", userID).Find(&albums)
	//查询出错返回
	if result.Error != nil {
		log.Printf("用户 %d 查询相册列表失败", userID)
		return nil, result.Error
	}
	log.Printf("用户 %d 共查询到 %d 个相册", userID, len(albums))
	return albums, result.Error
}

// QueryAlbumIDByUserIdAndAlbumName 根据用户ID和相册Name在albums表中查询该相册id
func QueryAlbumIDByUserIdAndAlbumName(db *gorm.DB, userID int, albumName string) int {
	var album Album
	result := db.Where("name = ? And user_id = ?", albumName, userID).First(&album)
	//查询出错返回
	if result.Error != nil {
		log.Printf("用户 %d 查询相册名为 %s 的相册ID失败", userID, albumName)
		return -1
	}
	log.Printf("用户 %d 查询相册名为 %s 的相册ID为：%d", userID, albumName, album.ID)
	return int(album.ID)
}

// UpdateAlbumNameByAlbumID 根据相册ID对albums表中的相册进行重命名
func UpdateAlbumNameByAlbumID(db *gorm.DB, albumID int, albumNewName string) error {
	if err := db.Model(&Album{}).Where("id = ?", albumID).Update("name", albumNewName).Error; err != nil {
		log.Printf("相册 %d 重命名为 %s 失败", albumID, albumNewName)
		return err
	}
	log.Printf("相册 %d 重命名为 %s", albumID, albumNewName)
	return nil
}

// DeleteAlbumByAlbumID 根据相册ID将albums表中的相册进行删除
func DeleteAlbumByAlbumID(db *gorm.DB, albumID int) error {
	//首先需要把items表中与album表关联的album_id字段中对应的albumID改成1，即默认相册，然后才是删除albums表中id为albumID的相册
	// 开启数据库事务
	tx := db.Begin()
	// 步骤1: 将与要删除的相册关联的照片的 album_id 更新为默认相册
	if err := tx.Model(&Item{}).Where("album_id = ?", albumID).Update("album_id", 1).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 步骤2: 删除 albums 表中的相册记录
	if err := tx.Where("id = ?", albumID).Delete(&Album{}).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
