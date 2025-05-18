package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type VerifyUserPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Token string `json:"token" validate:"required,min=5,max=200"`
}

func (app *application) verifyHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "verify")

	defer span.End()

	// get the parameters from
	var payload VerifyUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user verify payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating verify payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	span.SetAttributes(
		attribute.String("email", payload.Email),
	)
	user, err := app.store.Users().GetByEmailOrId(parentTraceCtx, &store.User{Email: payload.Email})
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to locate user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	token, err := app.store.Tokens().GetOne(parentTraceCtx, payload.Token, user.ID, store.VerificationTokenType)
	app.logger.Info(token, "TOKEN OOOO >>>>>>>>>>>>>>>>>>", user)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Trouble locating token", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if token.Value != payload.Token || token.ExpiresAt.Before(time.Now()) {
		err := errors.New("invalid token")
		app.logger.WithContext(parentTraceCtx).Error("Error validating verify payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.store.Tokens().Remove(parentTraceCtx, token)
	// unable to delete token
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to delete token", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Users().Verify(parentTraceCtx, nil, user)
	// unable to find user
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("verification error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "User verified successfully!", user)

}
