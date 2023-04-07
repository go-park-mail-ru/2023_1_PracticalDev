package postgres

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

type postgresRepository struct {
	db        *sql.DB
	log       log.Logger
	s3Service images.Service
}

func NewPostgresRepository(db *sql.DB, s3Service images.Service, log log.Logger) profile.Repository {
	return &postgresRepository{db, log, s3Service}
}

func (rep *postgresRepository) FullUpdate(params *profile.FullUpdateParams,
	image *models.Image) (profile.Profile, error) {
	const fullUpdateCmd = `UPDATE users
							SET username = $1::VARCHAR,
							name = $2::VARCHAR,
							profile_image = $3::VARCHAR,
							website_url = $4::VARCHAR
							WHERE id = $5
							RETURNING username, name, profile_image, website_url;`

	url, err := rep.s3Service.UploadImage(image)
	if err != nil {
		return profile.Profile{}, err
	}

	row := rep.db.QueryRow(fullUpdateCmd,
		params.Username,
		params.Name,
		url,
		params.WebsiteUrl,
		params.Id,
	)

	var prof profile.Profile
	var profileImage, websiteUrl sql.NullString
	err = row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	if err != nil {
		err = profile.ErrDb
	}
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}

func (rep *postgresRepository) PartialUpdate(params *profile.PartialUpdateParams) (profile.Profile, error) {
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
	var prof profile.Profile
	var profileImage, websiteUrl sql.NullString
	err := row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}
