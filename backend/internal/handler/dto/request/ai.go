package request

type CreateConversationRequest struct {
	Title string `json:"title"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}
