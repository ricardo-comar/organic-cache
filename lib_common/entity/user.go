package entity

import (
	"encoding/json"
	"strconv"
	"time"
)

type UserEntity struct {
	UserId string `dynamodbav:"user_id" json:"user_id"`
	TTL    string `dynamodbav:"ttl" json:"ttl"`
}

func NewUserEntity(body string) (*UserEntity, error) {

	user := &UserEntity{}
	json.Unmarshal([]byte(body), user)

	user.TTL = strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10)

	return user, nil
}
