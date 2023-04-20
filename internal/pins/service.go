package pins

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
)

type Pin struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MediaSource string `json:"media_source"`
	NumLikes    int    `json:"n_likes"`
	Liked       bool   `json:"liked"`
	Author      int    `json:"author_id"`
}

type Service interface {
	Create(params *CreateParams) (models.Pin, error)
	Get(id, userId int) (Pin, error)
	ListByAuthor(authorId, userId, page, limit int) ([]Pin, error)
	List(userId, page, limit int) ([]Pin, error)
	FullUpdate(params *FullUpdateParams) (models.Pin, error)
	Delete(id int) error

	CheckWriteAccess(userId, pinId string) (bool, error)
	CheckReadAccess(userId, pinId string) (bool, error)
}
