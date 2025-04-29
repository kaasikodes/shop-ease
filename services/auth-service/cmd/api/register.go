package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"github.com/kaasikodes/shop-ease/shared/kafka"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"go.opentelemetry.io/otel/codes"
)

type RegisterUserPayload struct {
	Email    string              `json:"email" validate:"required,email,max=255"`
	Name     string              `json:"name" validate:"required,max=100"`
	Password string              `json:"password" validate:"required,min=5,max=17"`
	RoleId   store.DefaultRoleID `json:"roleId" validate:"required"`
	Vendor   *VendorPayload      `json:"vendorInformation" validate:"-"`
}
type VendorPayload struct {
	Store struct {
		Name    string
		Address string
		Contact struct {
			phone string
			email string
		}
	}
	// TODO: Ensure all PII (personally identifiable info) is hashed
	Account struct {
		Bank       string
		AccNo      string
		SwiiftCode string
	}
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	registerTraceCtx, span := app.trace.Start(r.Context(), "register")

	defer span.End()

	// get the parameters from
	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(registerTraceCtx).Error("Error reading user registration payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(registerTraceCtx).Error("Error validating registration payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(registerTraceCtx, time.Second*5)
	defer cancel()
	span.SetAttributes(
		attribute.String("email", payload.Email),
		attribute.String("name", payload.Name),
		attribute.String("role", fmt.Sprintf("%v", payload.RoleId)),
	)
	// kafka test
	app.logger.Info("kafka lets go ...")
	err := kafka.PublishMessage(payload.Name, payload.Email)
	if err != nil {
		app.logger.Error(err, "Kafka error ....")
	}
	switch payload.RoleId {
	case store.CustomerID:
		app.logger.Info("New Customer Registeration initiated ...")
		user, verificationToken, err := app.registerCustomer(ctx, payload)
		if err != nil {
			app.logger.WithContext(registerTraceCtx).Error("Error registering customer", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			if err == store.ErrDuplicateEmail {

				app.badRequestResponse(w, r, err)
				return

			}
			app.internalServerError(w, r, err)
			return
		}
		if verificationToken != nil {

			vCtx := trace.ContextWithSpan(context.Background(), trace.SpanFromContext(registerTraceCtx)) // enures tracing id is maintained across services but ensure that the context is independent of the initial request context
			defer cancel()
			go func(ctx context.Context) {

				tokenCtx, span := app.trace.Start(ctx, "sending verification token")
				defer span.End()

				app.logger.WithContext(ctx).Info("Interaction with notification service begins  ....", *verificationToken)
				n, err := app.notificationService.Send(tokenCtx, &notification.NotificationRequest{
					Email:   user.Email,
					Title:   "Customer Verification",
					Content: fmt.Sprintf("This is your verification token %s", *verificationToken),
				},
				)
				if err != nil {
					app.logger.WithContext(ctx).Error("Error interacting with the notification service", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())

				} else {

					app.logger.Info("Success interacting with the notification service", n)
				}

			}(vCtx)
		}
		app.jsonResponse(w, http.StatusCreated, "Customer account created successfully, please check email for a verification link!", user)
		app.logger.WithContext(registerTraceCtx).Info("New Customer Registeration was a success ...", user.Email)

		return
	case store.VendorID:
		registerVendor()
		app.jsonResponse(w, http.StatusCreated, "Vendor account created successfully, please check email for a verification link!", nil)
		return
	default:
		app.badRequestResponse(w, r, errors.New("please select a valid role id"))
		return

	}

}
func (app *application) registerCustomer(ctx context.Context, payload RegisterUserPayload) (*store.User, *string, error) {
	_, span := app.trace.Start(ctx, "register-customer")
	defer span.End()
	// create user acc
	user := &store.User{
		Name:  payload.Name,
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		return nil, nil, err
	}
	plainToken := uuid.New().String() //TODO: hash and save token, not just the plain token
	err := app.store.Users().CreateWithVerificationToken(ctx, user, plainToken, ExpiresAtVerificationToken)

	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}

	// communicate with the notification service via grpc to send a verification email
	return user, &plainToken, nil

}
func registerVendor() {

}

// customer path - enter email, name, and password. Gets a verification mail via the notification service. The user clicks on the link on the verification mail(will contain a token) - token is valid, and then verified, and then account is verified & customer role is active or the token is invalid and an error message is sent to the user. Ensure validation shows proper error messages and path and all. All these events will be logged and how can this be used for data analysis, what role as well can prometheus and grafana play in this.
// vendor path
// a bit more info, and will have to interact with the subscription service first
