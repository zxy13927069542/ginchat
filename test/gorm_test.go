package test

import (
	"fmt"
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:066311@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "郑小燕"
	user.HeartbeatTime.Scan(time.Now())

	db.Create(user)

	// Read
	db.Debug().First(&user, "name = ?", "郑小燕")
	fmt.Println(user)

	// Update - 将 product 的 price 更新为 200
	db.Model(&user).Update("phone", "13927069542")
	// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&user)
}
