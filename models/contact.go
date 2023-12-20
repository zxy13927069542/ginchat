package models

import "gorm.io/gorm"

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerID  uint //	谁的关系信息
	TargetID uint //	对应的谁
	Type     int  //	对应类型
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
