package models

import (
	"errors"
	"gorm.io/gorm"
)

// GroupBasic 群组信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerID uint
	Icon    string
	Desc    string
	Type    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}


type GroupBasicModel struct {
	db *gorm.DB
}

func NewGroupBasicModel() *GroupBasicModel {
	return &GroupBasicModel{
		db: db,
	}
}

//	CreateGroup 新建群聊
func (m *GroupBasicModel) CreateGroup(group GroupBasic) error {
	if group.Name == "" || group.OwnerID <= 0 {
		return errors.New("群名或群主ID不能为空")
	}

	//	群主ID是否有效
	if NewUserBasicModel().FindUserByID(group.OwnerID).Name == "" {
		return errors.New("无效的群主ID")
	}

	//	群名是否重复
	var count int64
	m.db.Model(&GroupBasic{}).Where("name = ?", group.Name).Count(&count)
	if count >= 1 {
		return errors.New("群名重复")
	}

	return m.db.Create(&group).Error
}

//	FindGroupById 根据ID查询群聊
func (m *GroupBasicModel) FindGroupById(groupId uint) GroupBasic {
	var group GroupBasic
	m.db.Where("id = ?", groupId).First(&group)
	return group
}


