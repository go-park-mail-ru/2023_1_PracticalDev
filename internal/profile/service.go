package profile

import "errors"

var (
	ErrTooShortUsername = errors.New("username must be at least 4 characters")
	ErrTooLongUsername  = errors.New("username must be no more than 30 characters")
	ErrEmptyName        = errors.New("name must not be empty")
	ErrTooLongName      = errors.New("name must be no more than 60 characters")
)

type Service interface {
	FullUpdate(params *fullUpdateParams) (profile, error)
	PartialUpdate(params *partialUpdateParams) (profile, error)
}

func NewService(rep Repository) Service {
	return &service{rep}
}

type service struct {
	rep Repository
}

func validateUsername(username string) error {
	if len(username) < 4 {
		return ErrTooShortUsername
	} else if len(username) > 30 {
		return ErrTooLongUsername
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 1 {
		return ErrEmptyName
	} else if len(name) > 60 {
		return ErrTooLongName
	}
	return nil
}

func (serv *service) FullUpdate(params *fullUpdateParams) (profile, error) {
	if err := validateUsername(params.Username); err != nil {
		return profile{}, err
	} else if err = validateName(params.Name); err != nil {
		return profile{}, err
	}
	return serv.rep.FullUpdate(params)
}

func (serv *service) PartialUpdate(params *partialUpdateParams) (profile, error) {
	if params.UpdateUsername {
		if err := validateUsername(params.Username); err != nil {
			return profile{}, err
		}
	} else if params.UpdateName {
		if err := validateName(params.Name); err != nil {
			return profile{}, err
		}
	}
	return serv.rep.PartialUpdate(params)
}
