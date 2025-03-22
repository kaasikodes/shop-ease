package main

import "net/http"

func (app *application) verifyHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok__ verify",
		"environment": app.config.env,
		"version":     version,
	}
	if err := app.jsonResponse(w, http.StatusOK, "Health status retrieved successfully!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}
