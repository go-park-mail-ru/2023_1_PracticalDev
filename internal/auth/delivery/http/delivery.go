package http

import (
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/utils"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/tokens"
	mw "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/middleware"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
)

func RegisterHandlers(mux *httprouter.Router, logger *zap.Logger, serv auth.Service, token *tokens.HashToken, m *mw.HttpMetricsMiddleware) {
	del := delivery{serv, logger, token}
	mux.POST("/auth/login", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(del.Authenticate, logger), logger), logger))
	mux.DELETE("/auth/logout", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(del.Logout, logger), logger), logger))
	mux.POST("/auth/signup", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(del.Register, logger), logger), logger))
	mux.GET("/auth/me", mw.HandleLogger(mw.ErrorHandler(m.MetricsMiddleware(del.CheckAuth, logger), logger), logger))
}

type delivery struct {
	serv  auth.Service
	log   *zap.Logger
	token *tokens.HashToken
}

func parseSessionCookie(c *http.Cookie) (string, string, error) {
	tmp := strings.Split(c.Value, "$")
	if len(tmp) != 2 {
		return "", "", pkgErrors.ErrBadSessionCookie
	}

	return tmp[0], c.Value, nil
}

func createSessionCookie(s auth.SessionParams) *http.Cookie {
	return &http.Cookie{
		Name:     "JSESSIONID",
		Value:    s.Token,
		Expires:  time.Now().Add(s.LivingTime),
		HttpOnly: true,
		Path:     "/",
	}
}

func createCsrfTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 5),
		HttpOnly: false,
		Path:     "/",
	}
}

func (del *delivery) Authenticate(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	body, err := utils.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var request loginRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrParseJson, err.Error())
	}

	user, session, err := del.serv.Authenticate(request.Email, request.Password)
	if err != nil {
		return err
	}

	sessionCookie := createSessionCookie(session)
	http.SetCookie(w, sessionCookie)

	token, err := del.token.Create(&tokens.SessionParams{Token: session.Token}, time.Now().Add(session.LivingTime).Unix())
	if err != nil {
		del.log.Error("csrf token creation error:", zap.Error(err))
		return pkgErrors.ErrCreateCsrfToken
	}

	csrfCookie := createCsrfTokenCookie(token)
	http.SetCookie(w, csrfCookie)

	data, err := user.MarshalJSON()
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

func (del *delivery) CheckAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return pkgErrors.ErrBadSessionCookie
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return pkgErrors.ErrBadRequest
	}

	user, err := del.serv.CheckAuth(userId, sessionId)
	if err != nil {
		return err
	}

	data, err := user.MarshalJSON()
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

func (del *delivery) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	sessionCookie, err := r.Cookie("JSESSIONID")
	if err != nil {
		return pkgErrors.ErrUserNotFound
	}

	userId, sessionId, err := parseSessionCookie(sessionCookie)
	if err != nil {
		return pkgErrors.ErrBadRequest
	}

	if err = del.serv.DeleteSession(userId, sessionId); err != nil {
		return pkgErrors.ErrUnauthorized
	}

	newCookie := &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, newCookie)
	return pkgErrors.ErrNoContent
}

func (del *delivery) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	body, err := utils.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var request registerRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		return errors.Wrap(pkgErrors.ErrParseJson, err.Error())
	}

	params := auth.RegisterParams{
		Username: request.Username,
		Email:    request.Email,
		Name:     request.Name,
		Password: request.Password,
	}
	user, sessionParams, err := del.serv.Register(&params)
	if err != nil {
		return err
	}

	cookie := createSessionCookie(sessionParams)
	http.SetCookie(w, cookie)

	token, err := del.token.Create(&tokens.SessionParams{Token: sessionParams.Token}, time.Now().Add(sessionParams.LivingTime).Unix())
	if err != nil {
		del.log.Error("csrf token creation error:", zap.Error(err))
		return pkgErrors.ErrCreateCsrfToken
	}

	csrfCookie := createCsrfTokenCookie(token)
	http.SetCookie(w, csrfCookie)

	data, err := user.MarshalJSON()
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
