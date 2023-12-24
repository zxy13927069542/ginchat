package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromID   int64  //	发送者
	TargetID int64  //	接收者
	Type     int    //	消息类型 1-私信
	Media    int    //	消息类型 1-文字 图片 音频
	Content  string //	内容
	Pic      string //	图片
	Url      string
	Desc     string
	Amount   int //	其他数字统计

}

func (m *Message) TableName() string {
	return "message"
}
