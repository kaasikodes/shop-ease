package main

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

func (app *application) retriveAuthAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := app.trace.Start(ctx, "retrieving authenticated user details")

	user, ok := getUserFromContext(ctx)
	if !ok {
		err := errors.New("unable to retrieve user")
		app.logger.WithContext(ctx).Error("Retrieving user from context", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return

	}
	if err := app.jsonResponse(w, http.StatusOK, "Authenticated user retrieved successfully!", user); err != nil {
		app.internalServerError(w, r, err)
	}

}
