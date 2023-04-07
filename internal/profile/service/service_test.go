package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile/mocks"
)

func TestGetProfileByUser(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		id      int
		profile profile.Profile
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().GetProfileByUser(3).Return(profile.Profile{
					Username:     "un1",
					Name:         "n1",
					ProfileImage: "pi1",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			id:      3,
			profile: profile.Profile{Username: "un1", Name: "n1", ProfileImage: "pi1", WebsiteUrl: "wu1"},
			err:     nil,
		},
		"board not found": {
			prepare: func(f *fields) {
				f.repo.EXPECT().GetProfileByUser(3).Return(profile.Profile{}, profile.ErrProfileNotFound)
			},
			id:      3,
			profile: profile.Profile{},
			err:     profile.ErrProfileNotFound,
		},
		"negative user id param": {
			prepare: func(f *fields) {
				f.repo.EXPECT().GetProfileByUser(-1).Return(profile.Profile{}, profile.ErrProfileNotFound)
			},
			id:      -1,
			profile: profile.Profile{},
			err:     profile.ErrProfileNotFound,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewProfileService(f.repo)

			prof, err := serv.GetProfileByUser(test.id)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
		})
	}
}

func TestFullUpdate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  profile.FullUpdateParams
		profile profile.Profile
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().FullUpdate(&profile.FullUpdateParams{
					Id:           3,
					Username:     "username1",
					Name:         "n1",
					ProfileImage: models.Image{},
					WebsiteUrl:   "wu1",
				}).Return(profile.Profile{
					Username:     "username1",
					Name:         "n1",
					ProfileImage: "pi_url",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params: profile.FullUpdateParams{
				Id:           3,
				Username:     "username1",
				Name:         "n1",
				ProfileImage: models.Image{},
				WebsiteUrl:   "wu1",
			},
			profile: profile.Profile{Username: "username1", Name: "n1", ProfileImage: "pi_url", WebsiteUrl: "wu1"},
			err:     nil,
		},
		"too short username": {
			prepare: func(f *fields) {},
			params: profile.FullUpdateParams{
				Id:           3,
				Username:     "un1",
				Name:         "n1",
				ProfileImage: models.Image{},
				WebsiteUrl:   "wu1",
			},
			profile: profile.Profile{},
			err:     profile.ErrTooShortUsername,
		},
		"too long username": {
			prepare: func(f *fields) {},
			params: profile.FullUpdateParams{
				Id:           3,
				Username:     "username12_123456789_123456789_",
				Name:         "n1",
				ProfileImage: models.Image{},
				WebsiteUrl:   "wu1",
			},
			profile: profile.Profile{},
			err:     profile.ErrTooLongUsername,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewProfileService(f.repo)

			prof, err := serv.FullUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
		})
	}
}

func TestPartialUpdate(t *testing.T) {
	type fields struct {
		repo *mocks.MockRepository
	}

	type testCase struct {
		prepare func(f *fields)
		params  profile.PartialUpdateParams
		profile profile.Profile
		err     error
	}

	tests := map[string]testCase{
		"usual": {
			prepare: func(f *fields) {
				f.repo.EXPECT().PartialUpdate(&profile.PartialUpdateParams{
					Id:                 3,
					Username:           "username1",
					UpdateUsername:     true,
					Name:               "n1",
					UpdateName:         true,
					ProfileImage:       models.Image{},
					UpdateProfileImage: true,
					WebsiteUrl:         "wu1",
					UpdateWebsiteUrl:   true,
				}).Return(profile.Profile{
					Username:     "username1",
					Name:         "n1",
					ProfileImage: "pi_url",
					WebsiteUrl:   "wu1",
				}, nil)
			},
			params: profile.PartialUpdateParams{
				Id:                 3,
				Username:           "username1",
				UpdateUsername:     true,
				Name:               "n1",
				UpdateName:         true,
				ProfileImage:       models.Image{},
				UpdateProfileImage: true,
				WebsiteUrl:         "wu1",
				UpdateWebsiteUrl:   true,
			},
			profile: profile.Profile{Username: "username1", Name: "n1", ProfileImage: "pi_url", WebsiteUrl: "wu1"},
			err:     nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{repo: mocks.NewMockRepository(ctrl)}
			if test.prepare != nil {
				test.prepare(&f)
			}

			serv := NewProfileService(f.repo)

			prof, err := serv.PartialUpdate(&test.params)
			if err != test.err {
				t.Errorf("\nExpected: %s\nGot: %s", test.err, err)
			}
			if prof != test.profile {
				t.Errorf("\nExpected: %v\nGot: %v", test.profile, prof)
			}
		})
	}
}
