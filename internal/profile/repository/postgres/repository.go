package postgres

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/images"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/profile"
)

type postgresRepository struct {
	db      *sql.DB
	log     log.Logger
	imgServ images.Service
}

func NewPostgresRepository(db *sql.DB, imgServ images.Service, log log.Logger) profile.Repository {
	return &postgresRepository{db, log, imgServ}
}

func (rep *postgresRepository) GetProfileByUser(userId int) (profile.Profile, error) {
	const getCmd = `SELECT username, name, profile_image, website_url 
					FROM users 
					WHERE id = $1;`

	row := rep.db.QueryRow(getCmd, userId)

	var prof profile.Profile
	var profileImage, websiteUrl sql.NullString
	err := row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			err = profile.ErrProfileNotFound
		} else {
			err = profile.ErrDb
		}
	}
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}

const usernameUniqueConstraint = "users_username_key"

const fullUpdateCmd = `UPDATE users
						SET username = $1::VARCHAR,
						name = $2::VARCHAR,
						profile_image = $3::VARCHAR,
						website_url = $4::VARCHAR
						WHERE id = $5
						RETURNING username, name, profile_image, website_url;`

func (rep *postgresRepository) FullUpdate(params *profile.FullUpdateParams) (profile.Profile, error) {
	url, err := rep.imgServ.UploadImage(&params.ProfileImage)
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
		if err.(*pq.Error).Constraint == usernameUniqueConstraint {
			err = profile.ErrUsernameAlreadyExists
		} else {
			err = profile.ErrDb
		}
	}
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}

const partialUpdateCmd = `UPDATE users
							SET username = CASE WHEN $1::BOOLEAN THEN $2::VARCHAR ELSE username END,
							name = CASE WHEN $3::BOOLEAN THEN $4::VARCHAR ELSE name END,
							profile_image = CASE WHEN $5::BOOLEAN THEN $6::VARCHAR ELSE profile_image END,
							website_url = CASE WHEN $7::BOOLEAN THEN $8::VARCHAR ELSE website_url END
							WHERE id = $9
							RETURNING username, name, profile_image, website_url;`

func (rep *postgresRepository) PartialUpdate(params *profile.PartialUpdateParams) (profile.Profile, error) {
	var url string
	var err error
	if params.UpdateProfileImage {
		url, err = rep.imgServ.UploadImage(&params.ProfileImage)
		if err != nil {
			return profile.Profile{}, profile.ErrS3Service
		}
	}

	row := rep.db.QueryRow(partialUpdateCmd,
		params.UpdateUsername,
		params.Username,
		params.UpdateName,
		params.Name,
		params.UpdateProfileImage,
		url,
		params.UpdateWebsiteUrl,
		params.WebsiteUrl,
		params.Id,
	)
	var prof profile.Profile
	var profileImage, websiteUrl sql.NullString
	err = row.Scan(&prof.Username, &prof.Name, &profileImage, &websiteUrl)
	if err != nil {
		if err.(*pq.Error).Constraint == usernameUniqueConstraint {
			err = profile.ErrUsernameAlreadyExists
		} else {
			err = profile.ErrDb
		}
	}
	prof.ProfileImage = profileImage.String
	prof.WebsiteUrl = websiteUrl.String
	return prof, err
}

const isUsernameAvailableCmd = `SELECT NOT EXISTS(SELECT id
												FROM users
												WHERE username = $1 AND id <> $2) AS available;`

func (rep *postgresRepository) IsUsernameAvailable(username string, userId int) (bool, error) {
	row := rep.db.QueryRow(isUsernameAvailableCmd, username, userId)

	var available bool
	err := row.Scan(&available)
	if err != nil {
		err = profile.ErrDb
	}
	return available, err
}
