package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ricardo-comar/organic-cache/model"
)

func CreateEntity(body string) (model.UserEntity, error) {

	message := model.UserEntity{}
	json.Unmarshal([]byte(body), &message)

	message.TTL = strconv.FormatInt(time.Now().Add(time.Minute).UnixNano(), 10)

	return message, nil
}
