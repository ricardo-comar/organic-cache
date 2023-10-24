package entity

type UserEntity struct {
	UserId string `dynamodbav:"user_id" json:"user_id"`
	TTL    string `dynamodbav:"ttl" json:"ttl"`
}
