package boards

import (
	"errors"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type CreateParams struct {
	Name        string
	Description string
	Privacy     string
	UserId      int
}

type PartialUpdateParams struct {
	Id                int
	Name              string
	UpdateName        bool
	Description       string
	UpdateDescription bool
	Privacy           string
	UpdatePrivacy     bool
}

type FullUpdateParams struct {
	Id          int
	Name        string
	Description string
	Privacy     string
}

var (
	ErrBoardNotFound = errors.New("board not found")
	ErrDb            = errors.New("db error")
)

type Repository interface {
	Create(params *CreateParams) (models.Board, error)
	List(userId int) ([]models.Board, error)
	Get(id int) (models.Board, error)
	FullUpdate(params *FullUpdateParams) (models.Board, error)
	PartialUpdate(params *PartialUpdateParams) (models.Board, error)
	Delete(id int) error

	CheckWriteAccess(userId, boardId string) (bool, error)
	CheckReadAccess(userId, boardId string) (bool, error)
}
