package model

type UserEntity struct {
	ID string `dynamodbav:"id" json:"id"`
	// TTL string `dynamodbav:"ttl" json:"ttl"`
}
