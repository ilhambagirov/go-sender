package message

import "time"

type Dto struct {
	Id          uint
	Content     string
	Phone       string
	IsSent      bool
	CreatedDate time.Time
}

type PayloadDto struct {
	Content string `json:"content"`
	To      string `json:"to"`
}

type ResponseDto struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId"`
}
