package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromID   uint   //	发送者
	TargetID uint   //	接收者
	Type     string //	消息类型 群聊 私聊 广播	1-私信
	Media    int    //	消息类型 文字 图片 音频
	Contend  string //	内容
	Pic      string //	图片
	Url      string
	Desc     string
	Amount   int //	其他数字统计

}

func (m *Message) TableName() string {
	return "message"
}
