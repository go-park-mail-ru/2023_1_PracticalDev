package service

import (
	pkgComments "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/comments/mocks"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	notificationsMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/notifications/mocks"
	pinsMocks "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		params            *pkgComments.CreateParams
		comment           *models.Comment
	}

	type testCase struct {
		prepare func(f *fields)
		params  *pkgComments.CreateParams
		comment models.Comment
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Create(f.params).Return(*f.comment, nil),
					f.pinsRepo.EXPECT().Get(f.params.PinID).Return(models.Pin{Author: 12}, nil).
						MinTimes(0).MaxTimes(1),
					f.notificationsServ.EXPECT().Create(12, gomock.Any(), gomock.Any()).Return(nil).
						MinTimes(0).MaxTimes(1))
			},
			params:  &pkgComments.CreateParams{AuthorID: 27, PinID: 21, Text: "Good pin!"},
			comment: models.Comment{ID: 2, AuthorID: 27, PinID: 21, Text: "Good pin!"},
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				params:            test.params,
				comment:           &test.comment,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo)
			chat, err := serv.Create(test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if chat != test.comment {
				t.Errorf("\nExpected: %v\nGot: %v", test.comment, chat)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	type fields struct {
		repo              *mocks.MockRepository
		notificationsServ *notificationsMocks.MockService
		pinsRepo          *pinsMocks.MockRepository
		pinID             int
		comments          []models.Comment
	}

	type testCase struct {
		prepare  func(f *fields)
		pinID    int
		comments []models.Comment
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.pinID).Return(f.comments, nil)
			},
			pinID: 21,
			comments: []models.Comment{
				{ID: 2, AuthorID: 27, PinID: 21, Text: "Good pin!"},
				{ID: 3, AuthorID: 28, PinID: 21, Text: "Yeah!"},
				{ID: 4, AuthorID: 27, PinID: 21, Text: "Fantastic!"},
			},
			err: nil,
		},
		"no comments": {
			prepare: func(f *fields) {
				f.repo.EXPECT().List(f.pinID).Return(f.comments, nil)
			},
			pinID:    21,
			comments: []models.Comment{},
			err:      nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo:              mocks.NewMockRepository(ctrl),
				notificationsServ: notificationsMocks.NewMockService(ctrl),
				pinsRepo:          pinsMocks.NewMockRepository(ctrl),
				pinID:             test.pinID,
				comments:          test.comments,
			}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewService(f.repo, f.notificationsServ, f.pinsRepo)
			messages, err := serv.List(test.pinID)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if !reflect.DeepEqual(messages, test.comments) {
				t.Errorf("\nExpected: %v\nGot: %v", test.comments, messages)
			}
		})
	}
}
