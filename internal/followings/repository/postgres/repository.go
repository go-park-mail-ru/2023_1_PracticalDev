package postgres

import (
	"database/sql"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/followings"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type repository struct {
	db  *sql.DB
	log log.Logger
}

func NewRepository(db *sql.DB, log log.Logger) followings.Repository {
	return &repository{db, log}
}

const createCmd = `INSERT INTO followings (follower_id, followee_id)
				   VALUES ($1, $2);`

func (repo *repository) Create(followerId, followeeId int) error {
	_, err := repo.db.Exec(createCmd, followerId, followeeId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return nil
}

const deleteCmd = `DELETE FROM followings 
					WHERE follower_id = $1 AND followee_id = $2;`

func (repo *repository) Delete(followerId, followeeId int) error {
	_, err := repo.db.Exec(deleteCmd, followerId, followeeId)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	return nil
}

const getFolloweesCmd = `SELECT u.id, u.username, u.name, u.profile_image, u.website_url
							FROM users u
         							JOIN followings f ON u.id = f.followee_id
							WHERE f.follower_id = $1;`

func (repo *repository) GetFollowees(userId int) ([]followings.Followee, error) {
	rows, err := repo.db.Query(getFolloweesCmd, userId)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	followees := []followings.Followee{}
	followee := followings.Followee{}
	var profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&followee.Id, &followee.Username, &followee.Name, &profileImage, &websiteUrl)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		followee.ProfileImage = profileImage.String
		followee.WebsiteUrl = websiteUrl.String
		followees = append(followees, followee)
	}

	return followees, nil
}

const getFollowersCmd = `SELECT u.id, u.username, u.name, u.profile_image, u.website_url
							FROM users u
         							JOIN followings f ON u.id = f.follower_id
							WHERE f.followee_id = $1;`

func (repo *repository) GetFollowers(userId int) ([]followings.Follower, error) {
	rows, err := repo.db.Query(getFollowersCmd, userId)
	if err != nil {
		return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	followees := []followings.Follower{}
	followee := followings.Follower{}
	var profileImage, websiteUrl sql.NullString

	for rows.Next() {
		err = rows.Scan(&followee.Id, &followee.Username, &followee.Name, &profileImage, &websiteUrl)
		if err != nil {
			return nil, errors.Wrap(pkgErrors.ErrDb, err.Error())
		}

		followee.ProfileImage = profileImage.String
		followee.WebsiteUrl = websiteUrl.String
		followees = append(followees, followee)
	}

	return followees, nil
}

const userExistsCmd = `SELECT EXISTS(SELECT id
										FROM users
										WHERE id = $1) AS exists;`

func (repo *repository) UserExists(userId int) (bool, error) {
	row := repo.db.QueryRow(userExistsCmd, userId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return exists, nil
}

const followingExistsCmd = `SELECT EXISTS(SELECT follower_id
										FROM followings
										WHERE follower_id = $1 AND followee_id = $2) AS exists;`

func (repo *repository) FollowingExists(followerId, followeeId int) (bool, error) {
	row := repo.db.QueryRow(followingExistsCmd, followerId, followeeId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}
	return exists, nil
}
