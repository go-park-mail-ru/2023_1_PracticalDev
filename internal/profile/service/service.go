package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

type profileService struct {
	rep profile.Repository
}

func NewProfileService(rep profile.Repository) profile.Service {
	return &profileService{rep}
}

func validateUsername(username string) error {
	if len(username) < 4 {
		return profile.ErrTooShortUsername
	} else if len(username) > 30 {
		return profile.ErrTooLongUsername
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 1 {
		return profile.ErrEmptyName
	} else if len(name) > 60 {
		return profile.ErrTooLongName
	}
	return nil
}

func (serv *profileService) FullUpdate(params *profile.FullUpdateParams) (profile.Profile, error) {
	if err := validateUsername(params.Username); err != nil {
		return profile.Profile{}, err
	} else if err = validateName(params.Name); err != nil {
		return profile.Profile{}, err
	}
	return serv.rep.FullUpdate(params)
}

func (serv *profileService) PartialUpdate(params *profile.PartialUpdateParams) (profile.Profile, error) {
	if params.UpdateUsername {
		if err := validateUsername(params.Username); err != nil {
			return profile.Profile{}, err
		}
	} else if params.UpdateName {
		if err := validateName(params.Name); err != nil {
			return profile.Profile{}, err
		}
	}
	return serv.rep.PartialUpdate(params)
}
