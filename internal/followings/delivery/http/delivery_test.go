package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings/mocks"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDelivery_Follow(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  httprouter.Params
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Follow(3, 2).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "2"},
				{Key: "user-id", Value: "3"},
			},
			err: pkgErrors.ErrNoContent,
		},
		"invalid user id param": {
			prepare: nil,
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "a"},
			},
			err: pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare: nil,
			params:  []httprouter.Param{{Key: "id", Value: "3"}},
			err:     pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}
			defer logger.Sync()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodPost, "/users/3/following", nil)
			rec := httptest.NewRecorder()
			err = del.Follow(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
		})
	}
}

func TestDelivery_Unfollow(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare func(f *fields)
		params  httprouter.Params
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().Unfollow(2, 3).Return(nil)
			},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "2"},
			},
			err: pkgErrors.ErrNoContent,
		},
		"invalid user id param": {
			prepare: func(f *fields) {},
			params: []httprouter.Param{
				{Key: "id", Value: "3"},
				{Key: "user-id", Value: "a"},
			},
			err: pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare: func(f *fields) {},
			params:  []httprouter.Param{{Key: "id", Value: "3"}},
			err:     pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}
			defer logger.Sync()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodDelete, "/users/3/following", nil)
			rec := httptest.NewRecorder()
			err = del.Unfollow(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.err, err)
			}
		})
	}
}

func TestDelivery_GetFollowers(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetFollowers(12).Return([]followings.Follower{
					{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
					{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
					{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"followers":[{"id":2,"username":"vasua","name":"Vasya","profile_image":"vasya.jpg","website_url":"vasya.com"},{"id":3,"username":"kolya","name":"Kolya","profile_image":"kolya.jpg","website_url":"kolya.com"},{"id":4,"username":"sasha","name":"Sasha","profile_image":"sasha.jpg","website_url":"sasha.com"}]}`,
			err:      nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetFollowers(12).Return([]followings.Follower{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"followers":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}
			defer logger.Sync()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodGet, "/users/12/followers", nil)
			rec := httptest.NewRecorder()
			err = del.GetFollowers(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.response, string(body))
			}
		})
	}
}

func TestDelivery_GetFollowees(t *testing.T) {
	type fields struct {
		serv *mocks.MockService
	}

	type testCase struct {
		prepare  func(f *fields)
		params   httprouter.Params
		response string
		err      error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetFollowees(12).Return([]followings.Followee{
					{Id: 2, Username: "vasua", Name: "Vasya", ProfileImage: "vasya.jpg", WebsiteUrl: "vasya.com"},
					{Id: 3, Username: "kolya", Name: "Kolya", ProfileImage: "kolya.jpg", WebsiteUrl: "kolya.com"},
					{Id: 4, Username: "sasha", Name: "Sasha", ProfileImage: "sasha.jpg", WebsiteUrl: "sasha.com"},
				}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"followees":[{"id":2,"username":"vasua","name":"Vasya","profile_image":"vasya.jpg","website_url":"vasya.com"},{"id":3,"username":"kolya","name":"Kolya","profile_image":"kolya.jpg","website_url":"kolya.com"},{"id":4,"username":"sasha","name":"Sasha","profile_image":"sasha.jpg","website_url":"sasha.com"}]}`,
			err:      nil,
		},
		"no likes": {
			prepare: func(f *fields) {
				f.serv.EXPECT().GetFollowees(12).Return([]followings.Followee{}, nil)
			},
			params:   []httprouter.Param{{Key: "id", Value: "12"}},
			response: `{"followees":[]}`,
			err:      nil,
		},
		"invalid user id param": {
			prepare:  nil,
			params:   []httprouter.Param{{Key: "id", Value: "a"}},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
		"missing user id param": {
			prepare:  nil,
			params:   []httprouter.Param{},
			response: ``,
			err:      pkgErrors.ErrInvalidUserIdParam,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatalf("can't create logger: %s", err)
			}
			defer logger.Sync()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{serv: mocks.NewMockService(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			del := delivery{serv: f.serv, log: logger}
			req := httptest.NewRequest(http.MethodGet, "/users/12/followers", nil)
			rec := httptest.NewRecorder()
			err = del.GetFollowees(rec, req, test.params)
			if !errors.Is(err, test.err) {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			body, _ := io.ReadAll(rec.Body)
			if strings.Trim(string(body), "\n") != test.response {
				t.Errorf("\nExpected:\n%s\nGot:\n%s", test.response, string(body))
			}
		})
	}
}
