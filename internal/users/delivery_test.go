package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGet(t *testing.T) {
	res, _ := GetUser(10)
	assert.Equal(t, "Got user with id:10", res)
}
