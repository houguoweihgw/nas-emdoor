package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"nas-backend/entity"
	"time"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func InitDB() {
	// 连接数据库
	db, err := gorm.Open("mysql", "root:admin123@tcp(nas-mysql:3306)/nas_data?charset=utf8&parseTime=True&loc=Local")

	maxAttempts := 5 // 60 seconds (12 * 5 seconds)
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open("mysql", "root:admin123@tcp(nas-mysql:3306)/nas_data?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Printf("Attempt %d: Failed to connect to database. Retrying...\n", i)
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println("Connected to the database!")
			break
		}
	}

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// 设置数据库连接池等参数
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// 将数据库连接赋值给全局变量
	DB = db
	//// 注册关闭数据库连接的函数
	//defer func() {
	//	if err := DB.Close(); err != nil {
	//		panic("Failed to close database: " + err.Error())
	//	}
	//}()
}

func DBAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Item{})
}

func InitRedisClient() (err error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "nas-redis:6379",
		Password: "",
		DB:       0,
	})
	_, err = RDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
