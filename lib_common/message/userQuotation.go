package message

type UserPricesMessage struct {
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
}
