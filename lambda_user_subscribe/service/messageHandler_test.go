package service

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Body struct {
	ID string `json:"id"`
}

func TestIntMinTableDriven(t *testing.T) {

	body := Body{ID: "ABC"}
	bodyStr, _ := json.Marshal(body)

	entity, err := CreateEntity(string(bodyStr))

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, entity, "Unexpected nil entity")
	assert.NotNil(t, entity.UserId, "Unexpected nil id in entity")
	assert.NotNil(t, entity.TTL, "Unexpected nil ttl in entity")

	assert.Greater(t, entity.TTL, strconv.FormatInt(time.Now().UnixNano(), 10), "TTL must be greater than now")

}
