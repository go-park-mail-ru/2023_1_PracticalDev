package profile

type Service interface {
	GetProfileByUser(userId int) (Profile, error)
	FullUpdate(params *FullUpdateParams) (Profile, error)
	PartialUpdate(params *PartialUpdateParams) (Profile, error)
}
