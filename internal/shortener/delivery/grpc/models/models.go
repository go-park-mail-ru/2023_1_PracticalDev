package models

import shortener "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/shortener/delivery/grpc/proto"

func NewProtoStringMessage(msg string) *shortener.StringMessage {
	return &shortener.StringMessage{
		URL: msg,
	}
}

func NewStringMessage(msg *shortener.StringMessage) string {
	return msg.GetURL()
}
