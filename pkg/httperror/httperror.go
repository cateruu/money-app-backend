package httperror

import (
	"net/http"

	"github.com/cateruu/money-app-backend/pkg/json"
	"github.com/cateruu/money-app-backend/pkg/logger"
)

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	logger.Log.Error(err.Error(), "method", method, "uri", uri)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := json.Envelope{"error": message}

	err := json.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not procces your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid credentials provided"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func FailedAuthroizationResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "you do not have permission to access this resource"
	errorResponse(w, r, http.StatusUnauthorized, message)
}
