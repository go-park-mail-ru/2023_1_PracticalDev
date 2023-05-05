package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	hasherPkg "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type repository struct {
	db  *sql.DB
	rdb *redis.Client
	ctx context.Context
	log *zap.Logger
}

func NewRepository(db *sql.DB, rdb *redis.Client, ctx context.Context, log *zap.Logger) auth.Repository {
	return &repository{db, rdb, ctx, log}
}

func scanUser(user *models.User, row *sql.Row) error {
	var profileImage, websiteUrl sql.NullString
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.HashedPassword, &user.Name, &profileImage, &websiteUrl, &user.AccountType)
	user.WebsiteUrl = websiteUrl.String
	user.ProfileImage = profileImage.String
	return err
}

const authCommand = "SELECT * FROM users WHERE email = $1"

func (rep *repository) Authenticate(email, password string) (models.User, error) {
	const fnAuthenticate = "Authenticate"

	row := rep.db.QueryRow(authCommand, email)
	user := models.User{}
	hasher := hasherPkg.NewHasher()

	err := scanUser(&user, row)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, errors.Wrap(pkgErrors.ErrUserNotFound,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnAuthenticate,
				Query:  authCommand,
				Params: []any{email},
				Err:    err,
			}.Error())
	}

	if err != nil {
		return models.User{}, errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnAuthenticate,
				Query:  authCommand,
				Params: []any{email},
				Err:    err,
			}.Error())
	}

	if err = hasher.CompareHashAndPassword(user.HashedPassword, password); err != nil {
		return models.User{}, errors.Wrapf(pkgErrors.ErrWrongLoginOrPassword,
			"%s: error [%s]", fnAuthenticate, err)
	}

	return user, nil
}

func (rep *repository) SetSession(sessionId string, session *models.Session, expiration time.Duration) error {
	tmp, _ := json.Marshal(session)

	err := rep.rdb.HSet(rep.ctx, strconv.Itoa(session.UserId), sessionId, tmp).Err()
	if err != nil {
		rep.rdb.Expire(rep.ctx, strconv.Itoa(session.UserId), expiration)
	}

	return err
}

func (rep *repository) CheckAuth(userId, sessionId string) (models.User, error) {
	err := rep.rdb.HGet(rep.ctx, userId, sessionId).Err()
	user := models.User{}

	if err != nil {
		return user, err
	}

	row := rep.db.QueryRow("SELECT * FROM users WHERE id = $1", userId)
	err = scanUser(&user, row)

	return user, err
}

func (rep *repository) DeleteSession(userId, sessionId string) error {
	if err := rep.rdb.HGet(rep.ctx, userId, sessionId).Err(); err != nil {
		return err
	}
	rep.rdb.HDel(rep.ctx, userId, sessionId)
	return nil
}

func (rep *repository) Register(user *models.User) error {
	const fnRegister = "Register"
	const checkUserExistsCmd = "SELECT email FROM users WHERE email = $1"

	row := rep.db.QueryRow(checkUserExistsCmd, user.Email)
	tmp := ""
	err := row.Scan(&tmp)
	if err == nil {
		return errors.Wrap(pkgErrors.ErrUserAlreadyExists,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnRegister,
				Query:  checkUserExistsCmd,
				Params: []any{user.Email},
				Err:    err,
			}.Error())
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnRegister,
				Query:  checkUserExistsCmd,
				Params: []any{user.Email},
				Err:    err,
			}.Error())
	}

	const insertCommand = `INSERT INTO users (username, name, email, hashed_password, account_type, profile_image, website_url)
							VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = rep.db.Exec(insertCommand, user.Username, user.Name, user.Email, user.HashedPassword, user.AccountType, user.ProfileImage, user.WebsiteUrl)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrDb,
			pkgErrors.ErrRepositoryQuery{
				Func:   fnRegister,
				Query:  insertCommand,
				Params: []any{user.Username, user.Name, user.Email, user.HashedPassword, user.AccountType, user.ProfileImage, user.WebsiteUrl},
				Err:    err,
			}.Error())
	}

	return nil
}
