package postgres

import (
	"database/sql"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
)

func NewRepository(db *sql.DB, log *zap.Logger) users.Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log *zap.Logger
}

func (rep *repository) Get(id int) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE id = $1"
	var profileImage, websiteUrl sql.NullString
	rows, err := rep.db.Query(authCommand, id)
	if err != nil {
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	user := models.User{}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword, &user.Name, &profileImage,
		&websiteUrl, &user.AccountType)
	if err != nil {
		return models.User{}, errors.Wrap(pkgErrors.ErrDb, err.Error())
	}

	user.ProfileImage = profileImage.String
	user.WebsiteUrl = websiteUrl.String
	return user, err
}
