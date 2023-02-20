package users

import (
	"errors"
	"strconv"
)

func GetUser(id int) (string, error) {
	if id < 1000 {
		return "Got user with id:" + strconv.Itoa(id), nil
	}
	return "", errors.New("404: User not found")
}
