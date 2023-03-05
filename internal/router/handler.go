package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error
