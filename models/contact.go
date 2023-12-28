package models

import (
	"gorm.io/gorm"
)

// Contact 人员关系,诸如好友之类的
type Contact struct {
	gorm.Model
	OwnerID  uint //	谁的关系信息
	TargetID uint //	对应的谁
	Type     int  //	对应类型 1-models.Friend 2-models.Group 3-待定
	Desc     string
}

const (
	Friend = 1
	Group  = 2
)

func (table *Contact) TableName() string {
	return "contact"
}

type ContactModel struct {
	db *gorm.DB
}

func NewContactModel() *ContactModel {
	return &ContactModel{
		db: db,
	}
}

// SearchFriend 查询当前用户的好友信息
func (m *ContactModel) SearchFriend(userId uint) ([]UserBasic, error) {
	var friends []Contact
	if err := m.db.Where("owner_id = ? and type = ?", userId, Friend).Find(&friends).Error; err != nil {
		return nil, err
	}

	var list []uint
	for _, v := range friends {
		list = append(list, v.TargetID)
	}
	return UserBasicM.ListByIds(list)
}

//	SearchGroupFriends 根据群ID搜索群友ID
func (m *ContactModel) SearchGroupFriends(groupId uint) []uint {
	var friends []Contact
	m.db.Where("target_id = ? and type = ?", groupId, Group).Find(&friends)
	var ids []uint
	for _, v := range friends {
		ids = append(ids, v.OwnerID)
	}
	return ids
}
