package main

import (
	"net/http"
)

func (app *application) getAllNotifications(w http.ResponseWriter, r *http.Request) {
	// TODO: Retrieve query params and validate and pass them to get as opposed to nil, also create default pagination values -> should exist in shared
	result, total, err := app.store.Notification().Get(r.Context(), nil, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	notifications := make([]any, len(result))
	for i, v := range result {
		notifications[i] = v
	}
	if err = app.jsonResponse(w, http.StatusOK, "Notifications retrieved successfully!", paginatedResponse{Result: notifications, Total: total}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
