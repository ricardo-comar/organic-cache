package model

type SocketRequest struct {
	Action  string        `json:"action"`
	Payload SocketPayload `json:"payload"`
}

type SocketPayload struct {
	Message SocketMessage `json:"message"`
}

type SocketMessage struct {
	UserId      string        `dynamodbav:"user_id" json:"user_id"`
	ProductList []ProductItem `dynamodbav:"products" json:"products"`
}

type MessageEntity struct {
	RequestId string `json:"requestId"`
	UserId    string `json:"userId"`
}

type SocketResponse struct {
	Message string `json:"message"`
}
