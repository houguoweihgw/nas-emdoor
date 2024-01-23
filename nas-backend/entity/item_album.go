package entity

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type ItemAlbum struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ItemID    uint      `json:"item_id"`
	AlbumID   uint      `json:"album_id"`
	DateAdded time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_added"`
}

// TableName 设置表名
func (ItemAlbum) TableName() string {
	return "item_album"
}

// GetItemsInAlbum 查询item_album表获得表albumID中的所有照片
func GetItemsInAlbum(db *gorm.DB, albumID uint) ([]Item, error) {
	var photos []Item
	if err := db.
		Joins("JOIN item_album ON items.id = item_album.item_id").
		Where("item_album.album_id = ?", albumID).
		Find(&photos).Error; err != nil {
		log.Printf("查询相册 %d 所有照片失败", albumID)
		return nil, err
	}

	for i := range photos {
		db.Where("item_metadata.id = ? ", photos[i].ItemMetadataItem.ItemMetadataID).Find(&photos[i].ItemMetadata)
	}

	log.Printf("查询相册 %d 共有 %d 张照片", albumID, len(photos))
	return photos, nil
}

// AddPhotoToAlbum 将指定照片添加到指定相册
func AddPhotoToAlbum(db *gorm.DB, photoID int, albumID int) error {
	// 创建 item_album 记录来关联相册和照片
	itemAlbum := ItemAlbum{
		ItemID:  uint(photoID),
		AlbumID: uint(albumID),
	}
	if err := db.Create(&itemAlbum).Error; err != nil {
		log.Printf("用户将照片 %d 添加到相册 %d 失败", photoID, albumID)
		return err
	}
	log.Printf("用户将照片 %d 添加到相册 %d 成功", photoID, albumID)
	return nil
}

func RemovePhotoFromAlbum(db *gorm.DB, photoID int, albumID int) error {
	// 构建查询条件，删除指定照片与相册的关联关系
	result := db.Where("item_id = ? AND album_id = ?", photoID, albumID).Delete(&ItemAlbum{})
	// 检查删除过程中是否发生错误
	if result.Error != nil {
		log.Printf("将照片 %d 从相册 %d 中删除失败", photoID, albumID)
		return result.Error
	}
	// 检查是否成功删除了记录
	if result.RowsAffected == 0 {
		log.Printf("照片 %d 不再相册 %d 中，无法删除", photoID, albumID)
		return errors.New("关联记录不存在，无法移出照片")
	}
	log.Printf("将照片 %d 从相册 %d 中删除成功", photoID, albumID)
	return nil
}

// DeleteItemAlbumRecords 删除与指定照片相关的 item_album 记录
func DeleteItemAlbumRecords(tx *gorm.DB, photoID int) error {
	if err := tx.Where("item_id = ?", photoID).Delete(&ItemAlbum{}).Error; err != nil {
		log.Printf("将照片 %d 从所有相册中删除失败", photoID)
		return err
	}
	log.Printf("将照片 %d 从所有相册中删除成功", photoID)
	return nil
}
