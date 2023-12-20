package models

import (
	log "github.com/pion/ion-log"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identity      string
	Salt          string
	ClientIP      string
	ClientPort    string
	LoginTime     NullTime
	HeartbeatTime NullTime
	LoginOutTime  NullTime
	//LogOutTime    uint64
	IsLogOut   bool
	DeviceInfo string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

type UserBasicModel struct {
	db *gorm.DB
}

func NewUserBasicModel() *UserBasicModel {
	return &UserBasicModel{
		db: db,
	}
}

// FindUserByName()
func (m *UserBasicModel) FindUserByName(name string) UserBasic {
	var user UserBasic
	m.db.Where("name = ?", name).First(&user)
	return user
}

// FindUserByPhone()
func (m *UserBasicModel) FindUserByPhone(phone string) UserBasic {
	var user UserBasic
	m.db.Where("phone = ?", phone).First(&user)
	return user
}

// FindUserByEmail()
func (m *UserBasicModel) FindUserByEmail(email string) UserBasic {
	var user UserBasic
	m.db.Where("email = ?", email).First(&user)
	return user
}

// Create()
func (m *UserBasicModel) Create(user UserBasic) *gorm.DB {
	return m.db.Create(&user)
}

// Delete() user.ID不能为空
func (m *UserBasicModel) Delete(user UserBasic) *gorm.DB {
	return m.db.Delete(&user)
}

// Get() user.ID不能为空
func (m *UserBasicModel) Get(user *UserBasic) *gorm.DB {
	return m.db.First(user)
}

// Update() user.ID不能为空
func (m *UserBasicModel) Update(user UserBasic) *gorm.DB {
	return m.db.Updates(&user)
}

// List()
func (m *UserBasicModel) List() ([]UserBasic, error) {
	var users []UserBasic

	if err := m.db.Find(&users).Error; err != nil {
		log.Errorf("List UserBasic failed! Err: [%v]", err)
		return nil, err
	}
	return users, nil
}
