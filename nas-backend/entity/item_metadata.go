package entity

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type ItemMetadata struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	UserID       uint      `json:"user_id"`
	ExposureTime string    `json:"exposure_time"` //曝光时间
	Aperture     float64   `json:"aperture"`      //光圈值
	ISO          int       `json:"iso"`           //ISO
	FocalLength  float64   `json:"focal_length"`  //焦距
	Latitude     float64   `json:"latitude"`      //纬度
	Longitude    float64   `json:"longitude"`     //经度
	Altitude     float64   `json:"altitude"`      //高度
	Make         string    `json:"make"`          //相机品牌
	Model        string    `json:"model"`         //相机型号
	DateTaken    time.Time `json:"date_taken"`    //拍摄日期
	FileSize     int64     `json:"file_size"`     //照片大小
	ImageWidth   int       `json:"image_width"`   //照片宽度
	ImageLength  int       `json:"image_length"`  //照片高度
	SceneTags    string    `json:"scene_tags"`    //场景标签
}

// TableName 设置表名
func (ItemMetadata) TableName() string {
	return "item_metadata"
}

// InsertItemMetadata 构造照片结构体并插入item_metadata表
func InsertItemMetadata(db *gorm.DB, metadata ItemMetadata) (uint, error) {
	result := db.Create(&metadata)
	err := result.Error
	if err != nil {
		log.Printf("插入item_metadata表失败")
		return metadata.ID, err
	}
	log.Printf("插入item_metadata表成功，id：%d ", metadata.ID)
	return metadata.ID, err
}

// UpdateItemMetadataSceneTags 更新item_metadata表中对应id的场景识别标签scene_tags
func UpdateItemMetadataSceneTags(db *gorm.DB, id int, tags string) error {
	// 执行更新操作
	if err := db.Model(&ItemMetadata{}).Where("id = ?", id).Update("scene_tags", tags).Error; err != nil {
		log.Printf("更新item_metadata表失败，id: %d scene_tags: %s", id, tags)
		return err
	}
	log.Printf("更新item_metadata表成功，id: %d scene_tags: %s", id, tags)
	return nil
}

// GetCoversForLabel 为标签label查询第一张照片作为封面
func GetCoversForLabel(db *gorm.DB, label string, userID int) (string, error) {
	var coverURL string
	var itemMetadata ItemMetadata
	err := db.Where("user_id = ? And scene_tags LIKE ?", userID, "%"+label+"%").First(&itemMetadata).Error
	if err == nil {
		// 根据itemMetadata.ID查询item_metadata_item表和items来获得照片的地址
		// 如果找到匹配的标签，继续查询封面照片的地址
		var itemMetadataItem ItemMetadataItem
		err := db.Where("item_metadata_id = ?", itemMetadata.ID).First(&itemMetadataItem).Error
		if err == nil {
			// 如果找到匹配的 item_metadata_item，继续查询照片信息
			var item Item
			err := db.Where("id = ?", itemMetadataItem.ItemID).Last(&item).Error
			if err == nil {
				coverURL = item.FilePath
				log.Printf("用户id %d 查询场景分类标签 %s 封面成功", userID, label)
				return coverURL, nil
			}
		}
		return coverURL, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果未找到匹配项，将 URL 设置为空字符串或其他默认值
		log.Printf("用户id %d 查询场景分类标签 %s 封面为空", userID, label)
		return coverURL, err
	} else {
		log.Printf("用户id %d 查询场景分类标签 %s 封面失败", userID, label)
		return coverURL, err
	}
}

// GetUserLabelPhotos 获取用户标签的照片
func GetUserLabelPhotos(db *gorm.DB, userID int, labelName string, page int, pageSize int) ([]Item, error) {
	var items []Item
	// 构建查询条件
	query := db.Table("items").
		Select("items.*").Preload("ItemMetadataItem").
		Joins("JOIN item_metadata_item ON items.id = item_metadata_item.item_id").
		Joins("JOIN item_metadata ON item_metadata_item.item_metadata_id = item_metadata.id").
		Where("items.user_id = ? AND items.status = ? AND item_metadata.scene_tags LIKE ?",
			userID, "active", "%"+labelName+"%")
	// 计算分页偏移
	offset := (page - 1) * pageSize
	// 执行查询并应用分页
	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		log.Printf("用户id %d 查询场景分类标签 %s 的照片为空", userID, labelName)
		return nil, err
	}
	for i := range items {
		db.Where("item_metadata.id = ? ", items[i].ItemMetadataItem.ItemMetadataID).Find(&items[i].ItemMetadata)
	}
	log.Printf("用户id %d 查询场景分类标签 %s 的照片成功", userID, labelName)
	return items, nil
}
