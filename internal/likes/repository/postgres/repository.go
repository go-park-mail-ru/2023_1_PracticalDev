package postgres

import (
	"database/sql"
	pkgLikes "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/likes"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type repository struct {
	db  *sql.DB
	log log.Logger
}

func NewRepository(db *sql.DB, log log.Logger) pkgLikes.Repository {
	return &repository{db, log}
}

const createCmd = `INSERT INTO pin_likes (pin_id, author_id)
				   VALUES ($1, $2);`

func (repo *repository) Create(pinId, authorId int) error {
	res, err := repo.db.Exec(createCmd, pinId, authorId)
	if err != nil {
		return pkgLikes.ErrDb
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected < 1 {
		return pkgLikes.ErrDb
	}
	return nil
}

const deleteCmd = `DELETE FROM pin_likes 
					WHERE pin_id = $1 AND author_id = $2;`

func (repo *repository) Delete(pinId, authorId int) error {
	res, err := repo.db.Exec(deleteCmd, pinId, authorId)
	if err != nil {
		return pkgLikes.ErrDb
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected < 1 {
		return pkgLikes.ErrDb
	}
	return nil
}

const listByAuthorCmd = `SELECT pin_id, author_id, created_at
						FROM pin_likes 
						WHERE author_id = $1
						ORDER BY created_at DESC;`

func (repo *repository) ListByAuthor(authorId int) ([]models.Like, error) {
	rows, err := repo.db.Query(listByAuthorCmd, authorId)
	if err != nil {
		return nil, pkgLikes.ErrDb
	}

	likes := []models.Like{}
	like := models.Like{}
	for rows.Next() {
		err = rows.Scan(&like.PinId, &like.AuthorId, &like.CreatedAt)
		if err != nil {
			return nil, pkgLikes.ErrDb
		}
		likes = append(likes, like)
	}
	return likes, nil
}

const listByPinCmd = `SELECT pin_id, author_id, created_at
					FROM pin_likes 
					WHERE pin_id = $1
					ORDER BY created_at DESC;`

func (repo *repository) ListByPin(pinId int) ([]models.Like, error) {
	rows, err := repo.db.Query(listByPinCmd, pinId)
	if err != nil {
		return nil, pkgLikes.ErrDb
	}

	likes := []models.Like{}
	like := models.Like{}
	for rows.Next() {
		err = rows.Scan(&like.PinId, &like.AuthorId, &like.CreatedAt)
		if err != nil {
			return nil, pkgLikes.ErrDb
		}
		likes = append(likes, like)
	}
	return likes, nil
}

const pinExistsCmd = `SELECT EXISTS(SELECT id
									FROM pins
									WHERE id = $1) AS exists;`

func (repo *repository) PinExists(pinId int) (bool, error) {
	row := repo.db.QueryRow(pinExistsCmd, pinId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, pkgLikes.ErrDb
	}
	return exists, nil
}

const userExistsCmd = `SELECT EXISTS(SELECT id
										FROM users
										WHERE id = $1) AS exists;`

func (repo *repository) UserExists(userId int) (bool, error) {
	row := repo.db.QueryRow(userExistsCmd, userId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, pkgLikes.ErrDb
	}
	return exists, nil
}

const likeExistsCmd = `SELECT EXISTS(SELECT pin_id
										FROM pin_likes
										WHERE pin_id = $1 AND author_id = $2) AS exists;`

func (repo *repository) LikeExists(pinId, authorId int) (bool, error) {
	row := repo.db.QueryRow(likeExistsCmd, pinId, authorId)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, pkgLikes.ErrDb
	}
	return exists, nil
}
