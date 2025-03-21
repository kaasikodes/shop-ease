package main

import "net/http"

func (app *application) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok____",
		"environment": app.config.env,
		"version":     version,
	}
	if err := app.jsonResponse(w, http.StatusOK, "Health status retrieved successfully!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}