package message

type CreateMessageRequest struct {
	Content string `json:"content"`
	Phone   string `json:"phone"`
}
