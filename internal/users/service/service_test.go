package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestService_Get(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
		id   int
		user models.User
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		user    models.User
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(f.user, nil)
			},
			id: 3,
			user: models.User{Id: 3, Username: "petya", Email: "petya@vk.com", Name: "Petya",
				ProfileImage: "petya.jpg", WebsiteUrl: "petya.com", AccountType: "personal"},
			err: nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(f.user, pkgErrors.ErrProfileNotFound)
			},
			id:   3,
			user: models.User{},
			err:  pkgErrors.ErrProfileNotFound,
		},
		"negative user id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().Get(f.id).Return(f.user, pkgErrors.ErrProfileNotFound)
			},
			id:   -1,
			user: models.User{},
			err:  pkgErrors.ErrProfileNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl), id: test.id, user: test.user}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo)
			user, err := serv.Get(test.id)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(user, test.user) {
				t.Errorf("\nExpected: %v\nGot: %v", test.user, user)
			}
		})
	}
}
