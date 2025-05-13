package main

import "net/http"

func (app *application) healthzHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Health route reached ...")
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.Env,
		"version":     version,
		"service":     "Notification",
	}
	if err := app.jsonResponse(w, http.StatusOK, "Health status retrieved successfully!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}
