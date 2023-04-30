package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	pkgChats "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/chats"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

const (
	chatsUrl = "/chats"
	chatUrl  = "/chats/:id"
)

func RegisterHandlers(mux *httprouter.Router, logger log.Logger, authorizer mw.Authorizer, serv pkgChats.Service) {
	del := delivery{serv, logger}

	mux.POST(chatsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.create)), logger), logger))
	mux.GET(chatsUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.listByUser)), logger), logger))
	mux.GET(chatUrl, mw.HandleLogger(mw.ErrorHandler(mw.Cors(authorizer(del.get)), logger), logger))
}

type delivery struct {
	serv pkgChats.Service
	log  log.Logger
}

func (del *delivery) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	decoder := json.NewDecoder(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			del.log.Error(err)
		}
	}()

	var request createRequest
	if err = decoder.Decode(&request); err != nil {
		return errors.Wrap(pkgErrors.ErrParseJson, err.Error())
	}

	params := pkgChats.CreateParams{User1ID: userID, User2ID: request.UserID}
	chat, err := del.serv.Create(&params)
	if err != nil {
		return err
	}

	response := newCreateResponse(&chat)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) listByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strUserID := p.ByName("user-id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidUserIdParam, err.Error())
	}

	chats, err := del.serv.ListByUser(userID)
	if err != nil {
		return err
	}

	response := newListResponse(chats)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	strID := p.ByName("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrInvalidChatIDParam, err.Error())
	}

	chat, err := del.serv.Get(id)
	if err != nil {
		return err
	}

	response := newGetResponse(&chat)
	data, err := json.Marshal(response)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrCreateResponse, err.Error())
	}
	return nil
}
