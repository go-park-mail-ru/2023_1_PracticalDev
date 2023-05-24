package service

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/boards"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgPins "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pins"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/constants"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

type service struct {
	pinServ pkgPins.Service
	repo    boards.Repository
}

func NewBoardsService(repo boards.Repository, pinServ pkgPins.Service) boards.Service {
	return &service{repo: repo, pinServ: pinServ}
}

func validateName(name string) error {
	if len(name) > constants.MaxBoardNameLen {
		return pkgErrors.ErrTooLongBoardName
	}
	return nil
}

func validateDescription(name string) error {
	if len(name) > constants.MaxBoardDescriptionLen {
		return pkgErrors.ErrTooLongBoardDescription
	}
	return nil
}

func validatePrivacy(privacy string) error {
	if privacy != "secret" && privacy != "public" {
		return pkgErrors.ErrInvalidPrivacy
	}
	return nil
}

func (serv *service) Create(params *boards.CreateParams) (models.Board, error) {
	if err := validatePrivacy(params.Privacy); err != nil {
		return models.Board{}, err
	} else if err = validateName(params.Name); err != nil {
		return models.Board{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Board{}, err
	}

	return serv.repo.Create(params)
}

func (serv *service) List(userId int) ([]models.Board, error) {
	return serv.repo.List(userId)
}

func (serv *service) Get(id int) (models.Board, error) {
	return serv.repo.Get(id)
}

func (serv *service) FullUpdate(params *boards.FullUpdateParams) (models.Board, error) {
	if err := validatePrivacy(params.Privacy); err != nil {
		return models.Board{}, err
	} else if err = validateName(params.Name); err != nil {
		return models.Board{}, err
	} else if err = validateDescription(params.Description); err != nil {
		return models.Board{}, err
	}

	return serv.repo.FullUpdate(params)
}

func (serv *service) PartialUpdate(params *boards.PartialUpdateParams) (models.Board, error) {
	if params.UpdateName {
		if err := validateName(params.Name); err != nil {
			return models.Board{}, err
		}
	}
	if params.UpdateDescription {
		if err := validateDescription(params.Description); err != nil {
			return models.Board{}, err
		}
	}
	if params.UpdatePrivacy {
		if err := validatePrivacy(params.Privacy); err != nil {
			return models.Board{}, err
		}
	}

	return serv.repo.PartialUpdate(params)
}

func (serv *service) Delete(id int) error {
	return serv.repo.Delete(id)
}

func (serv *service) AddPin(boardId, pinId int) error {
	exists, err := serv.repo.HasPin(boardId, pinId)
	if err != nil {
		return err
	}
	if exists {
		return pkgErrors.ErrPinAlreadyAdded
	}

	return serv.repo.AddPin(boardId, pinId)
}

func (serv *service) PinsList(userId, boardId int, page, limit int) ([]models.Pin, error) {
	pins, err := serv.repo.PinsList(boardId, page, limit)
	if err != nil {
		return pins, err
	}

	for i := range pins {
		err = serv.pinServ.SetLikedField(&pins[i], userId)
		if err != nil {
			return pins, err
		}
	}

	return pins, err
}

func (serv *service) RemovePin(boardId, pinId int) error {
	return serv.repo.RemovePin(boardId, pinId)
}

func (serv *service) CheckWriteAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckWriteAccess(userId, boardId)
}

func (serv *service) CheckReadAccess(userId, boardId string) (bool, error) {
	return serv.repo.CheckReadAccess(userId, boardId)
}
