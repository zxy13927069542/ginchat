package models

import "gorm.io/gorm"

// GroupBasic 群组信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerID string
	Icon    string
	Desc    string
	Type    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
