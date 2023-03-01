package middleware

// import (
// 	"encoding/json"
// 	"net/http"
// 	"github.com/julienschmidt/httprouter"
// )
//
// type response struct {
//   title string
//   status int
//   detail string
// }
//
//
// func Authorize(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, rep *auth.Repository) httprouter.Handle {
//   return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//   session_id, err := r.Cookie("JSESSIONID")
//   
//   if err != nil {
//     w.WriteHeader(http.StatusNotFound)
//     resp := response {
//       title: "Error",
//       status: http.StatusNotFound,
//       detail: "Cookie not found",
//     }
//     tmp, _:= json.Marshal(resp)
//     w.Write(tmp)
//   }
//
//   _, err := del.serv.CheckAuth(session_id) 
//
//   if err != nil {
//     return errors.New("user is not authorized")
//   }
//
//   j, _ := json.Marshal(user)
//   w.Write(j)
//   return nil
//   }
// }
