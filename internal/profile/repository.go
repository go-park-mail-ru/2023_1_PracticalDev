package profile

import (
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Profile struct {
	Username     string
	Name         string
	ProfileImage string
	WebsiteUrl   string
}

type FullUpdateParams struct {
	Id           int
	Username     string
	Name         string
	ProfileImage models.Image
	WebsiteUrl   string
}

type PartialUpdateParams struct {
	Id                 int
	Username           string
	UpdateUsername     bool
	Name               string
	UpdateName         bool
	ProfileImage       models.Image
	UpdateProfileImage bool
	WebsiteUrl         string
	UpdateWebsiteUrl   bool
}

var (
	ErrDb                    = errors.New("db error")
	ErrS3Service             = errors.New("s3 service error")
	ErrProfileNotFound       = errors.New("profile not found")
	ErrUsernameAlreadyExists = errors.New("user with such username already exists")
)

type Repository interface {
	GetProfileByUser(userId int) (Profile, error)
	FullUpdate(params *FullUpdateParams) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)

	IsUsernameAvailable(username string, userId int) (bool, error)
}
