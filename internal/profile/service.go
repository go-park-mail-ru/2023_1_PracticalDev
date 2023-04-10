package profile

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrTooShortUsername = errors.New("username must be at least 4 characters")
	ErrTooLongUsername  = errors.New("username must be no more than 30 characters")
	ErrEmptyName        = errors.New("name must not be empty")
	ErrTooLongName      = errors.New("name must be no more than 60 characters")
)

type ErrBadParams struct {
	Err error
}

func (e ErrBadParams) Error() string {
	return fmt.Sprintf("profile service: bad params: %s", e.Err)
}

type Service interface {
	GetProfileByUser(userId int) (Profile, error)
	FullUpdate(params *FullUpdateParams) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)
}
