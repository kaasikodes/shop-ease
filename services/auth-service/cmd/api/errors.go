package main

import "net/http"

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}
	writeJsonError(w, http.StatusInternalServerError,"The server encountered an unexpected issue", errors)

}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path, "error")
	errors := []string {}

	writeJsonError(w, http.StatusForbidden, "forbidden", errors)
}



func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}
	writeJsonError(w, http.StatusBadRequest, "Validation Error", errors)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}
	writeJsonError(w, http.StatusConflict, "Conflict", errors)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}

	writeJsonError(w, http.StatusNotFound, "not found", errors)
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}

	writeJsonError(w, http.StatusUnauthorized, "unauthorized", errors)
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errors := []string {}
	if !app.isProduction() {
		errors[0] = err.Error()
		
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJsonError(w, http.StatusUnauthorized, "unauthorized", errors)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	errors := []string {}

	w.Header().Set("Retry-After", retryAfter)

	writeJsonError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter, errors)
}