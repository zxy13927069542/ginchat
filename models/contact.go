package models

import (
	"errors"
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

//	AddFriend 添加好友
func (m *ContactModel) AddFriend(userId uint, targetName string) error {
	//	查询目标用户是否存在
	var userByName UserBasic
	if userByName = NewUserBasicModel().FindUserByName(targetName); userByName.Name != targetName {
		return errors.New("用户不存在")
	}

	if userId == userByName.ID {
		return errors.New("不能添加自己为好友")
	}

	//	查询是否已经是好友关系
	var count int64
	m.db.Model(&Contact{}).Where("owner_id = ? and target_id = ? and type = ?", userId, userByName.ID, Friend).
		Or("owner_id = ? and target_id = ? and type = ?", userByName.ID, userId, Friend).Count(&count)
	if count > 0 {
		return errors.New("好友关系已存在")
	}

	tx := m.db.Begin()
	//	发生panic()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//	创建两条好友关系
	var c1 Contact
	c1.OwnerID = userId
	c1.TargetID = userByName.ID
	c1.Type = Friend
	if err := m.db.Create(&c1).Error; err != nil {
		tx.Rollback()
		return err
	}
	var c2 Contact
	c2.OwnerID = userByName.ID
	c2.TargetID = userId
	c2.Type = Friend
	if err := m.db.Create(&c2).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
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

//	AddGroup 添加群聊
func (m *ContactModel) AddGroup(userId uint, group string) error {
	//	群聊是否存在
	g := NewGroupBasicModel().FindGroupByIdOrName(group)
	if g.Name == "" {
		return errors.New("群聊不存在")
	}

	//	是否已入群
	var count int64
	m.db.Model(&Contact{}).Where("owner_id = ? and target_id = ? and type = ?", userId, g.ID, Group).Count(&count)
	if count >= 1 {
		return errors.New("该用户已在群内")
	}

	var contact Contact
	contact.OwnerID = userId
	contact.TargetID = g.ID
	contact.Type = Group
	return m.db.Create(&contact).Error
}

//	LoadGroup 查询用户所在群聊
func (m *ContactModel) LoadGroup(ownerId uint) []GroupBasic {
	var contacts []Contact
	var groupIds []uint
	m.db.Where("owner_id = ? and type = ?", ownerId, Group).Find(&contacts)
	for _, v := range contacts {
		groupIds = append(groupIds, v.TargetID)
	}

	var groupList []GroupBasic
	m.db.Where("id in ?", groupIds).Find(&groupList)
	return groupList
}
