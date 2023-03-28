package profile

import (
	"database/sql"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
)

type profile struct {
	Username     string
	Name         string
	ProfileImage string
	WebsiteUrl   string
}

type fullUpdateParams struct {
	Id           int
	Username     string
	Name         string
	ProfileImage string
	WebsiteUrl   string
}

type partialUpdateParams struct {
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
	FullUpdate(params *fullUpdateParams) (profile, error)
	PartialUpdate(params *partialUpdateParams) (profile, error)
}

func NewRepository(db *sql.DB, log log.Logger) Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) FullUpdate(params *fullUpdateParams) (profile, error) {
	const fullUpdateCmd = `UPDATE users
								SET username = $1::VARCHAR,
    							name = $2::VARCHAR,
    							profile_image = $3::VARCHAR,
								website_url = $4::VARCHAR
								WHERE id = $5
								RETURNING username, name, profile_image, website_url;`

	row := rep.db.QueryRow(fullUpdateCmd,
		params.Username,
		params.Name,
		params.ProfileImage,
		params.WebsiteUrl,
		params.Id,
	)
	var prof profile
	var profileImage, websiteUrl sql.NullString
	err := row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}

func (rep *repository) PartialUpdate(params *partialUpdateParams) (profile, error) {
	const partialUpdateBoard = `UPDATE users
								SET username = CASE WHEN $1::BOOLEAN THEN $2::VARCHAR ELSE username END,
								name = CASE WHEN $3::BOOLEAN THEN $4::VARCHAR ELSE name END,
    							profile_image = CASE WHEN $5::BOOLEAN THEN $6::VARCHAR ELSE profile_image END,
    							website_url = CASE WHEN $7::BOOLEAN THEN $8::VARCHAR ELSE website_url END
								WHERE id = $9
								RETURNING username, name, profile_image, website_url;`

	row := rep.db.QueryRow(partialUpdateBoard,
		params.UpdateUsername,
		params.Username,
		params.UpdateName,
		params.Name,
		params.UpdateProfileImage,
		params.ProfileImage,
		params.UpdateWebsiteUrl,
		params.WebsiteUrl,
		params.Id,
	)
	var prof profile
	var profileImage, websiteUrl sql.NullString
	err := row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}
