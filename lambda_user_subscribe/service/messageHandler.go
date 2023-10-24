package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ricardo-comar/organic-cache/lib_common/entity"
)

func CreateEntity(body string) (entity.UserEntity, error) {

	message := entity.UserEntity{}
	json.Unmarshal([]byte(body), &message)

	message.TTL = strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10)

	return message, nil
}
