package elapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	ElApi struct {
		router *mux.Router
	}
	ApiFunc  func(http.ResponseWriter, *http.Request) error
	ApiError struct {
		Error string
	}
)

func NewElApi() *ElApi {
	route := mux.NewRouter()
	return &ElApi{
		router: route,
	}
}

func (ap *ElApi) GetRouter() *mux.Router {
	return ap.router
}

func (ap *ElApi) WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func (ap *ElApi) MakeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			ap.WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (ap *ElApi) PermissionDenied(w http.ResponseWriter) {
	ap.WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func (ap *ElApi) GetFromVars(r *http.Request, elvar string) string {
	return mux.Vars(r)[elvar]
}

func (ap *ElApi) Route(elroute string, elfunc ApiFunc, method string) {
	ap.router.HandleFunc(elroute, ap.MakeHTTPHandleFunc(elfunc)).Methods(method)
}
