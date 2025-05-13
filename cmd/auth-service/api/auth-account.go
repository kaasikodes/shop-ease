package main

import "net/http"

func (app *application) retriveAuthAccountHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok Auth",
		"environment": app.config.env,
		"version":     version,
	}
	if err := app.jsonResponse(w, http.StatusOK, "Health status retrieved successfully!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}
