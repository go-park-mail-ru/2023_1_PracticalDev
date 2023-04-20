package postgres

import (
	"database/sql"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
)

func NewRepository(db *sql.DB, log log.Logger) users.Repository {
	return &repository{db, log}
}

type repository struct {
	db  *sql.DB
	log log.Logger
}

func (rep *repository) Get(id int) (models.User, error) {
	authCommand := "SELECT * FROM users WHERE id = $1"
	var profileImage, websiteUrl sql.NullString
	rows, err := rep.db.Query(authCommand, id)

	if err != nil {
		return models.User{}, err
	}

	user := models.User{}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword, &user.Name, &profileImage,
		&websiteUrl, &user.AccountType)
	user.ProfileImage = profileImage.String
	user.WebsiteUrl = websiteUrl.String
	return user, err
}
