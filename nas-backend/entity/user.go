package entity

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// User 结构用于表示用户数据
type User struct {
	Id       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"not null;unique" json:"email"`
}

// ValidLoginInfo 根据用户名查询users表验证登陆信息，返回是否验证成功
func ValidLoginInfo(db *gorm.DB, user User) bool {
	var existingUser User
	// 利用用户名从数据库中查询用户
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		// 处理查询错误，这里可以根据具体需求进行处理
		return false
	}
	if user.Username != existingUser.Username {
		return false
	}
	//log.Printf("name:%s  user:%s", user.Username, existingUser)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// 验证密码
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(existingUser.Password))
	if err != nil {
		// 密码验证失败
		return false
	}
	// 密码验证成功
	return true
}

// GetUserIdByUsername 根据用户名查询users表，返回用户id
func GetUserIdByUsername(db *gorm.DB, username string) int {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return -1
	}
	log.Printf("用户 %s 查询id成功：%d", username, user.Id)
	return int(user.Id)
}

// UserRegister 用户注册
func UserRegister(db *gorm.DB, user User) error {
	result := db.Create(&user)
	err := result.Error
	if err != nil {
		log.Printf("用户：%s注册失败", user.Username)
		return err
	}
	log.Printf("users表插入成功, 用户名：%s, 用户邮箱：%s", user.Username, user.Email)
	return nil
}
