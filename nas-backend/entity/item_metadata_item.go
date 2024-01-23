package entity

import (
	"github.com/jinzhu/gorm"
	"log"
)

type ItemMetadataItem struct {
	ItemMetadataID uint `gorm:"foreignKey:ItemMetadataID;references:ID"`
	ItemID         uint `gorm:"foreignKey:ItemID;references:ID"`
}

// TableName 设置表名
func (ItemMetadataItem) TableName() string {
	return "item_metadata_item"
}

// CreateItemMetadataItem 插入item_metadata_item表
func CreateItemMetadataItem(db *gorm.DB, itemMetadataID, itemID uint) error {
	newItemMetadataItem := ItemMetadataItem{
		ItemMetadataID: itemMetadataID,
		ItemID:         itemID,
	}
	result := db.Create(&newItemMetadataItem)
	if result.Error != nil {
		log.Printf("itemMetadataID：%d ItemID:%d 插入item_metadata_item表失败\n", itemMetadataID, itemID)
		return result.Error
	}
	log.Printf("itemMetadataID：%d ItemID:%d 插入item_metadata_item表成功\n", itemMetadataID, itemID)
	return nil
}
