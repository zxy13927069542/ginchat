package models

import (
	"errors"
	log "github.com/pion/ion-log"
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

// CreateGroup 新建群聊
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

	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//	先建群
	if err := m.db.Create(&group).Error; err != nil {
		log.Errorf(">>CreateGroup() failed! Err: [%v]", err)
		tx.Rollback()
		return errors.New("创建失败")
	}

	//	将群主ID加入群聊
	var contact Contact
	contact.OwnerID = group.OwnerID
	contact.TargetID = group.ID
	contact.Type = Group
	if err := m.db.Create(&contact).Error; err != nil {
		log.Errorf(">>CreateGroup() failed! Err: [%v]", err)
		tx.Rollback()
		return errors.New("创建失败")
	}
	tx.Commit()
	return nil
}

// FindGroupById 根据ID查询群聊
func (m *GroupBasicModel) FindGroupByIdOrName(group string) GroupBasic {
	var g GroupBasic
	m.db.Where("id = ? or name = ?", group, group).First(&g)
	return g
}
