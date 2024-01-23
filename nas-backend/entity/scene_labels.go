package entity

import (
	"github.com/jinzhu/gorm"
	"log"
)

type SceneLabel struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	UserID     int    `gorm:"not null" json:"user_id"`
	LabelName  string `gorm:"type:varchar(255);not null" json:"label_name"`
	LabelCount int    `gorm:"default:0" json:"label_count"`
}

// TableName 设置表名
func (SceneLabel) TableName() string {
	return "scene_labels"
}

// UpdateOrInsertLabelCount 更新标签数量
func UpdateOrInsertLabelCount(db *gorm.DB, labelName string, userID int) error {
	// 首先尝试查找标签
	label := SceneLabel{}
	if err := db.Where("label_name = ? AND user_id = ?", labelName, userID).First(&label).Error; err == nil {
		// 标签已存在，递增计数
		label.LabelCount++
		db.Save(&label)
	} else {
		// 标签不存在，创建新记录
		newLabel := SceneLabel{
			LabelName:  labelName,
			LabelCount: 1,
			UserID:     userID,
		}
		db.Create(&newLabel)
	}
	return nil
}

// ParseSceneClassificationLabels 解析场景分类标签
func ParseSceneClassificationLabels(db *gorm.DB, userID int, tags []string) {
	for _, tag := range tags {
		err := UpdateOrInsertLabelCount(db, tag, userID)
		if err != nil {
			return
		}
	}
	log.Printf("标签 %v 更新完成\n", tags)
}

// GetSceneLabelsByUserID 获取用户标签
func GetSceneLabelsByUserID(db *gorm.DB, userID int) ([]string, error) {
	var labels []SceneLabel
	var result []string

	// 查询用户的标签，其中 LabelCount 大于零
	if err := db.Where("user_id = ? AND label_count > 0", userID).Find(&labels).Error; err != nil {
		log.Printf("用户 %d 查询所有场景分类标签失败: %v\n", userID, err)
		return nil, err
	}

	// 提取标签名称
	for _, label := range labels {
		result = append(result, label.LabelName)
	}
	log.Printf("用户 %d 查询所有场景分类标签成功: %v\n", userID, result)
	return result, nil
}

// GetSceneLabelsInfoByUserID 获取用户标签
func GetSceneLabelsInfoByUserID(db *gorm.DB, userID int) ([]SceneLabel, error) {
	var labels []SceneLabel
	// 查询用户的标签，其中 LabelCount 大于零
	if err := db.Where("user_id = ? AND label_count > 1", userID).
		Order("label_count desc, label_name").
		Find(&labels).Error; err != nil {
		log.Printf("用户 %d 查询所有场景分类标签失败: %v\n", userID, err)
		return nil, err
	}
	log.Printf("用户 %d 查询所有场景分类标签成功\n", userID)
	return labels, nil
}

// GetLabelTotalPhotosCount  获取用户标签的照片数量
func GetLabelTotalPhotosCount(db *gorm.DB, userID int, labelName string) (int, error) {
	var label SceneLabel
	// 查询用户的标签，其中 LabelCount 大于零
	if err := db.Where("user_id = ? AND label_name = ?", userID, labelName).Find(&label).Error; err != nil {
		log.Printf("用户 %d 查询场景分类标签 %s 照片数量失败: %v\n", userID, labelName, err)
		return 0, err
	}
	log.Printf("用户 %d 查询场景分类标签 %s 照片数量成功: %d\n", userID, labelName, label.LabelCount)
	return label.LabelCount, nil
}

// GetSearchLabelTotalPhotosCount 获取用户搜索标签的照片数量
func GetSearchLabelTotalPhotosCount(db *gorm.DB, userID int, labelName string) (int, error) {
	var searchCount int
	var label []SceneLabel
	// 查询用户的标签，其中 LabelCount 大于零
	if err := db.Where("user_id = ? AND label_name LIKE ?", userID, "%"+labelName+"%").Find(&label).Error; err != nil {
		log.Printf("用户 %d 查询场景分类标签 %s 照片数量失败: %v\n", userID, labelName, err)
		return 0, err
	}
	for _, l := range label {
		searchCount += l.LabelCount
	}
	log.Printf("用户 %d 查询场景分类标签 %s 照片数量成功: %d\n", userID, labelName, searchCount)
	return searchCount, nil
}
