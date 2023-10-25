package service

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ricardo-comar/organic-cache/lib_common/entity"
)

func CreateEntity(body string) (*entity.UserEntity, error) {

	user := &entity.UserEntity{}
	json.Unmarshal([]byte(body), user)

	user.TTL = strconv.FormatInt(time.Now().Add(time.Minute*5).UnixNano(), 10)

	return user, nil
}
