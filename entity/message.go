package entity

import "time"

type message struct {
	Room     string    `json:"room" binding:"required"`
	Author   string    `json:"author" binding:"required"`
	Text     string    `json:"text"   binding:"required"`
	CreateAt time.Time `json:"time"   binding:"required"`
}

func NewEmptyMessage() *message {
	return &message{}
}

func NewMessage(room string, author string, text string) *message {
	return &message{
		Room:     room,
		Author:   author,
		Text:     text,
		CreateAt: time.Now(),
	}
}
