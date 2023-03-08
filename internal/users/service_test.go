package users

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type serviceTestCase struct {
	id      int
	user    models.User
	IsError bool
}

func TestServiceGetUser(t *testing.T) {
	cases := []serviceTestCase{
		{
			id: 1,
			user: models.User{
				Id:             1,
				Username:       "geogreck",
				Email:          "geogreck@vk.com",
				HashedPassword: "$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS",
				Name:           "George",
				ProfileImage:   "",
				WebsiteUrl:     "",
				AccountType:    "personal",
			},
			IsError: false,
		},
		{
			id: 2,
			user: models.User{
				Id:             2,
				Username:       "kirill",
				Email:          "figma@vk.com",
				HashedPassword: "$2a$10$A4Ab/cuy/oLNvm4VxGoO/ezKL.fiew5e.eKTkUOWIVxoBh8XFO4lS",
				Name:           "Kirill",
				ProfileImage:   "",
				WebsiteUrl:     "",
				AccountType:    "personal",
			},
			IsError: false,
		},
		{
			id:      0,
			user:    models.User{},
			IsError: true,
		},
		{
			id:      100,
			user:    models.User{},
			IsError: true,
		},
		{
			id:      -1,
			user:    models.User{},
			IsError: true,
		},
	}

	logger := log.New()
	db, err := db.New(logger)
	if err != nil {
		os.Exit(1)
	}
	serv := NewService(NewRepository(db, logger))

	var user models.User
	for _, item := range cases {
		user, err = serv.GetUser(item.id)
		isError := err != nil

		assert.Equal(t, item.IsError, isError)
		assert.Equal(t, item.user, user)
	}
}
