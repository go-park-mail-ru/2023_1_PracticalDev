package middleware

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/router"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AccessService interface {
	CheckWriteAccess(userId, objectId string) (bool, error)
	CheckReadAccess(userId, objectId string) (bool, error)
}

type AccessChecker struct {
	serv AccessService
}

func NewAccessChecker(serv AccessService) AccessChecker {
	return AccessChecker{serv}
}

func (accessChecker *AccessChecker) WriteChecker(handler router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		userId := p.ByName("user-id")
		objectId := p.ByName("id")

		access, err := accessChecker.serv.CheckWriteAccess(userId, objectId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		if !access {
			w.WriteHeader(http.StatusForbidden)
			return nil
		}

		return handler(w, r, p)
	}
}

func (accessChecker *AccessChecker) ReadChecker(handler router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		userId := p.ByName("user-id")
		objectId := p.ByName("id")

		access, err := accessChecker.serv.CheckReadAccess(userId, objectId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		if !access {
			w.WriteHeader(http.StatusForbidden)
			return nil
		}

		return handler(w, r, p)
	}
}
