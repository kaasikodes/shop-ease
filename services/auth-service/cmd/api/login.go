package main

import (
	"errors"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=5,max=17"`
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "login")

	defer span.End()

	// get the parameters from
	var payload LoginUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user login payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating login payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	span.SetAttributes(
		attribute.String("email", payload.Email),
	)
	user, err := app.store.Users().GetByEmailOrId(parentTraceCtx, &store.User{
		Email: payload.Email,
	})
	// unable to find user
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Unable to locate user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	// is user verified
	if !user.IsVerified {
		err = errors.New("user is not verified")
		app.logger.WithContext(parentTraceCtx).Error("Verification Error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	// check if passwords match
	passwordsMatch := user.Password.Compare(payload.Password)

	if !passwordsMatch {
		err = errors.New("passwords do not match")
		app.logger.WithContext(parentTraceCtx).Error("Verification Error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return

	}

	app.jsonResponse(w, http.StatusOK, "User logged in successfully!", map[string]any{"user": user, "accessToken": ""})

}
