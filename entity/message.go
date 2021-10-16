package entity

import "time"

type message struct {
	Room     string    `json:"room" binding:"required"`
	UserName string    `json:"userName" binding:"required"`
	Text     string    `json:"text"   binding:"required"`
	CreateAt time.Time `json:"time"   binding:"required"`
}

func NewEmptyMessage() *message {
	return &message{}
}

func NewMessage(room string, userName string, text string) *message {
	return &message{
		Room:     room,
		UserName: userName,
		Text:     text,
		CreateAt: time.Now(),
	}
}
