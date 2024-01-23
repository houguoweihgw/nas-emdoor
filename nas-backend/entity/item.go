package entity

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"log"
	"time"
)

type Item struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	FilePath         string    `json:"-"`
	DateUploaded     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"dateUploaded"`
	UserID           uint      // 与 users 表的关联字段
	Status           string    `gorm:"default:'active'" json:"status"`
	Collected        bool      `json:"collected"`
	ItemMetadataItem ItemMetadataItem
	ItemMetadata     ItemMetadata
}

// TableName 使用 GORM 的约定指定关联关系
func (Item) TableName() string {
	return "items"
}

// GetAllPhotosByUserID 根据用户id查询items表，返回该用户的所有激活照片
func GetAllPhotosByUserID(db *gorm.DB, userid int) ([]Item, error) {
	var items []Item
	//查询该用户所有照片
	result := db.Where("user_id = ? And status =?", userid, "active").Find(&items)
	//查询出错返回
	if result.Error != nil {
		log.Printf("用户 %d 查询照片失败", userid)
		return nil, result.Error
	}
	log.Printf("用户 %d 共查询到 %d 张照片", userid, len(items))
	return items, result.Error
}

// GetTotalPhotosCountByUserID 根据用户id查询items表，返回该用户的激活照片数量
func GetTotalPhotosCountByUserID(db *gorm.DB, userID int) (int, error) {
	var count int
	result := db.Model(&Item{}).Where("user_id = ? AND status = ?", userID, "active").Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// GetPagePhotosByUserID 根据用户id查询items表，返回这一页该用户的所有激活照片（即从offset到offset+perPage）
func GetPagePhotosByUserID(db *gorm.DB, userID int, offset int, perPage int) ([]Item, error) {
	var items []Item
	result := db.Preload("ItemMetadataItem").
		Joins("INNER JOIN item_metadata_item ON items.id = item_metadata_item.item_id").
		Where("items.user_id = ? AND items.status = ?", userID, "active").
		Offset(offset).
		Limit(perPage).
		Find(&items)

	if result.Error != nil {
		log.Printf("用户 %d 查询照片失败", userID)
		return nil, result.Error
	}

	for i := range items {
		db.Where("item_metadata.id = ? ", items[i].ItemMetadataItem.ItemMetadataID).Find(&items[i].ItemMetadata)
	}

	log.Printf("用户 %d 共查询到 %d 张照片", userID, len(items))
	return items, result.Error

}

// GetWelcomePhotosByUserID 根据用户id查询items表，返回随机的六张激活照片
func GetWelcomePhotosByUserID(db *gorm.DB, userID int, count int) ([]Item, error) {
	var items []Item
	result := db.Where("user_id = ? AND status = ?", userID, "active").
		Order("RAND()").Limit(6).
		Find(&items)

	if result.Error != nil {
		log.Printf("用户 %d 查询欢迎照片失败", userID)
		return nil, result.Error
	}

	log.Printf("用户 %d 共查询到 %d 张欢迎照片", userID, len(items))
	return items, result.Error

}

// GetAllCollectedPhotosByUserID 根据用户id查询items表，返回该用户的所有激活照片
func GetAllCollectedPhotosByUserID(db *gorm.DB, userid int) ([]Item, error) {
	var items []Item
	//查询该用户所有收藏照片
	result := db.Preload("ItemMetadataItem").
		Joins("INNER JOIN item_metadata_item ON items.id = item_metadata_item.item_id").
		Where("user_id = ? And status = ? And collected = ?", userid, "active", true).
		Find(&items)
	//查询出错返回
	if result.Error != nil {
		log.Printf("用户 %d 查询收藏照片失败", userid)
		return nil, result.Error
	}
	for i := range items {
		db.Where("item_metadata.id = ? ", items[i].ItemMetadataItem.ItemMetadataID).Find(&items[i].ItemMetadata)
	}
	log.Printf("用户 %d 共查询到 %d 张已收藏照片", userid, len(items))
	return items, result.Error
}

// GetAllRecycledPhotosByUserID 根据用户id查询items表，返回该用户的所有回收站照片
func GetAllRecycledPhotosByUserID(db *gorm.DB, userid int) ([]Item, error) {
	var items []Item
	//查询该用户所有照片
	result := db.Preload("ItemMetadataItem").
		Joins("INNER JOIN item_metadata_item ON items.id = item_metadata_item.item_id").
		Where("user_id = ? And status =?", userid, "recycled").
		Find(&items)
	//查询出错返回
	if result.Error != nil {
		log.Printf("用户 %d 查询照片失败", userid)
		return nil, result.Error
	}
	for i := range items {
		db.Where("item_metadata.id = ? ", items[i].ItemMetadataItem.ItemMetadataID).Find(&items[i].ItemMetadata)
	}
	log.Printf("用户 %d 共查询到 %d 张照片", userid, len(items))
	return items, result.Error
}

// UpdatePhotoCollectedByPhotoID 根据照片id查询items表，更改collected状态
func UpdatePhotoCollectedByPhotoID(db *gorm.DB, photoId int) error {
	// 将items表中数据库中collected字段取反
	if err := db.Model(&Item{}).Where("id = ?", photoId).Update("collected", gorm.Expr("NOT collected")).Error; err != nil {
		log.Printf("更改照片 %d collected状态失败", photoId)
		return err
	}
	log.Printf("更改照片 %d collected状态成功", photoId)
	return nil
}

// RemovePhotoByPhotoID 根据照片id将items表中的照片状态更新为recycled(回收站)
func RemovePhotoByPhotoID(db *gorm.DB, photoID int) error {
	tx := db.Begin()
	// 开始事务
	if tx.Error != nil {
		return tx.Error
	}
	// 删除与该照片相关的所有 item_album 记录
	if err := DeleteItemAlbumRecords(tx, photoID); err != nil {
		tx.Rollback()
		log.Printf("清除照片 %d 相册关联失败", photoID)
		return err
	}
	// 更新 items 表中数据库中对应照片的状态改为“recycled”
	if err := tx.Model(&Item{}).Where("id = ?", photoID).Update("status", "recycled").Error; err != nil {
		tx.Rollback()
		log.Printf("清除照片 %d 失败", photoID)
		return err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("事务提交失败")
		return err
	}
	log.Printf("清除照片 %d 成功", photoID)
	return nil
}

// RecoverPhotoByPhotoID 根据照片id将items表中的照片状态更新为active
func RecoverPhotoByPhotoID(db *gorm.DB, photoID int) error {
	// 将items表中数据库中对应照片的状态改为“recycled”
	if err := db.Model(&Item{}).Where("id = ?", photoID).Update("status", "active").Error; err != nil {
		log.Printf("恢复照片 %d 失败", photoID)
		return err
	}
	log.Printf("恢复照片 %d 成功", photoID)
	return nil
}

// DeletePhotoByPhotoID 根据照片id将items表中的照片状态更新为deleted(彻底删除)
func DeletePhotoByPhotoID(db *gorm.DB, photoID int) error {
	// 将items表中数据库中对应照片的状态改为“deleted”
	if err := db.Model(&Item{}).Where("id = ?", photoID).Update("status", "deleted").Error; err != nil {
		log.Printf("删除照片 %d 失败", photoID)
		return err
	}
	log.Printf("删除照片 %d 成功", photoID)
	return nil
}

// InsertPhoto 构造照片结构体并插入items表
func InsertPhoto(db *gorm.DB, title string, description string, filePath string, userID int) (uint, error) {
	item := Item{
		Title:       title,
		Description: description,
		FilePath:    filePath,
		UserID:      uint(userID),
	}
	result := db.Create(&item)
	err := result.Error
	if err != nil {
		log.Printf("items表插入失败, 用户ID：%d", userID)
		return item.ID, err
	}
	log.Printf("items表插入成功, 用户ID：%d", userID)
	return item.ID, nil
}

// InsertPhotoWithMetadata 将照片插入到items表,以及将照片的元数据插入到item_metadata表
func InsertPhotoWithMetadata(db *gorm.DB, item *Item, itemMeta *ItemMetadata) error {
	//由于items表和item_metadata表的设计问题,需要上传照片时先给item_metadata_id一个临时值,
	//等item_metadata表添加完成后,最后将items表的item_metadata_id更新
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 第一部,插入items表
	if err := tx.Create(&item).Error; err != nil {
		log.Printf("items表插入失败,事务回滚")
		tx.Rollback()
		return err
	}
	log.Printf("items表插入成功, 照片ID：%d", item.ID)
	// 第二步,插入item_metadata表
	if err := tx.Create(&itemMeta).Error; err != nil {
		log.Printf("item_metadata表插入失败,事务回滚")
		tx.Rollback()
		return err
	}
	log.Printf("item_metadata表插入成功, 元数据ID：%d", itemMeta.ID)
	// 第三步更新中间表item_metadata_item，关联items和item_metadata
	newItemMetadataItem := ItemMetadataItem{
		ItemMetadataID: itemMeta.ID,
		ItemID:         item.ID,
	}
	if err := tx.Create(&newItemMetadataItem).Error; err != nil {
		log.Printf("item_metadata_item表插入失败,事务回滚")
		tx.Rollback()
		return err
	}
	log.Printf("item_metadata_item表插入成功, ItemMetadataID：%d  ItemID:%d", newItemMetadataItem.ItemMetadataID, newItemMetadataItem.ItemID)
	return tx.Commit().Error
}
