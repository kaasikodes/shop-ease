package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/oauth/provider"
	store_base "github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	store "github.com/kaasikodes/shop-ease/services/auth-service/internal/store/sql-store"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) githubOauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "Github Oauth Authorization Login")

	defer span.End()
	githubProvider, ok := app.oauthProviderRegistry[provider.OauthProviderTypeGithub]
	if !ok {
		err := errors.New("unregistered oauth provider")
		app.logger.WithContext(parentTraceCtx).Error("Oauth provider error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.WithContext(parentTraceCtx).Error("Oauth provider loggin in ...")
	_roleId := r.URL.Query().Get("roleId")
	roleId, err := strconv.Atoi(_roleId)
	if err != nil {
		err := errors.Join(err, errors.New("invalid role id"))
		app.logger.WithContext(parentTraceCtx).Error("Invalid roleId", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, errors.New("please provide a valid roleId"))
		return
	}
	githubProvider.Login(w, r, &provider.LoginOption{RoleId: store_base.DefaultRoleID(roleId)})
}
func (app *application) githubOauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "Github Oauth Authorization Callback")

	defer span.End()
	provider, ok := app.oauthProviderRegistry[provider.OauthProviderTypeGithub]
	if !ok {
		err := errors.New("unregistered oauth provider")
		app.logger.WithContext(parentTraceCtx).Error("Oauth provider error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.logger.WithContext(parentTraceCtx).Error(err, "oauth provider error")
		app.badRequestResponse(w, r, err)
		return
	}

	info, err := provider.Callback(w, r)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Oauth provider error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.logger.WithContext(parentTraceCtx).Error(err, "oauth provider callback error")
		app.badRequestResponse(w, r, err)
		return
	}
	app.logger.Info(info, "LOgggin in ....", *info, &info)
	// path 1: user exists: create jwt token and attach to response for user
	user, err := app.store.Users().GetByEmailOrId(parentTraceCtx, &store.User{Email: info.Email})
	if err != nil && !errors.Is(err, store_base.ErrNoUserFound) {
		app.logger.WithContext(parentTraceCtx).Error("User store error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.logger.WithContext(parentTraceCtx).Error(err, "issue locating user")
		app.notFoundResponse(w, r, err)
		return
	}
	// path 2: user does not exists, check the attached role and then move appropriately
	if err != nil && errors.Is(err, store_base.ErrNoUserFound) {
		switch info.RoleId {
		case store.CustomerID:

			if user, _, err := app.registerCustomer(parentTraceCtx, RegisterUserPayload{Email: info.Email, Name: info.Name}, true); err == nil {
				accessToken, err := app.jwt.CreateToken(strconv.Itoa(user.ID), user.Email, AccessTokenDuration)
				if err != nil {
					app.logger.WithContext(parentTraceCtx).Error("Jwt token err", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					app.badRequestResponse(w, r, err)
					return

				}

				app.jsonResponse(w, http.StatusOK, "User logged in successfully", LoginResponse{
					User:        *user,
					AccessToken: accessToken,
				})
				return
			} else {
				app.logger.WithContext(parentTraceCtx).Error("Customer Registeration Error", err)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				app.logger.WithContext(parentTraceCtx).Error(err, "Customer Registeration Error")
				app.badRequestResponse(w, r, err)

			}
		default:
			app.badRequestResponse(w, r, errors.New("please select a valid role id"))
			return
		}

	}
	accessToken, err := app.jwt.CreateToken(strconv.Itoa(user.ID), user.Email, AccessTokenDuration)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Jwt token err", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return

	}

	app.jsonResponse(w, http.StatusOK, "User logged in successfully", LoginResponse{
		User:        *user,
		AccessToken: accessToken,
	})
}
