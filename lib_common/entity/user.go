package entity

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserEntity struct {
	UserId string `dynamodbav:"user_id" json:"user_id"`
	TTL    string `dynamodbav:"ttl" json:"ttl"`
}

func NewUserEntity(body string) (*UserEntity, error) {

	user := &UserEntity{}
	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Invalid json content: %v / error: %v", body, err)
		return nil, err
	}
	if len(user.UserId) == 0 {
		log.Printf("Empty user_id")
		return nil, errors.New("user_id_empty")
	}

	user.TTL = strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10)

	return user, nil
}
