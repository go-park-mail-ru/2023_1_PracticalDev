package profile

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

type Repository interface {
	FullUpdate(params *FullUpdateParams) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)
}
