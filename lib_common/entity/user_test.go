package entity

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {

	entity, err := NewUserEntity("{ \"user_id\":\"ABC\"}")

	assert.Nil(t, err, "Unexpected error")
	assert.NotNil(t, entity, "Unexpected nil entity")
	assert.NotNil(t, entity.UserId, "Unexpected nil id in entity")
	assert.NotNil(t, entity.TTL, "Unexpected nil ttl in entity")

	assert.Greater(t, entity.TTL, strconv.FormatInt(time.Now().UnixNano(), 10), "TTL must be greater than now")

}

func TestEmpty(t *testing.T) {

	entity, err := NewUserEntity("")

	assert.NotNil(t, err, "Unexpected nil error")
	assert.Nil(t, entity, "Unexpected nil entity")
}

func TestUnexpectedFields(t *testing.T) {

	entity, err := NewUserEntity("{ \"id\":\"ABC\"}")

	assert.NotNil(t, err, "Unexpected nil error")
	assert.Nil(t, entity, "Unexpected nil entity")
}

func TestEmptyUserID(t *testing.T) {

	entity, err := NewUserEntity("{ \"ttl\":\"ABC\"}")

	assert.NotNil(t, err, "Unexpected nil error")
	assert.Nil(t, entity, "Unexpected nil entity")
}
