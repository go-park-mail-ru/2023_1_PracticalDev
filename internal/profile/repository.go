package profile

import (
	"errors"

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
	ProfileImage string
	WebsiteUrl   string
}

type PartialUpdateParams struct {
	Id                 int
	Username           string
	UpdateUsername     bool
	Name               string
	UpdateName         bool
	ProfileImage       string
	UpdateProfileImage bool
	WebsiteUrl         string
	UpdateWebsiteUrl   bool
}

var (
	ErrDb = errors.New("db error")
)

type Repository interface {
	FullUpdate(params *FullUpdateParams, image *models.Image) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)
}
