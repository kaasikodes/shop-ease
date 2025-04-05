package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
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

	// get the parameters from
	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*1)
	defer cancel()
	switch payload.RoleId {
	case store.CustomerID:
		user, verificationToken, err := app.registerCustomer(ctx, payload)
		if err != nil {
			if err == store.ErrDuplicateEmail {
				app.badRequestResponse(w, r, err)
				return

			}
			app.internalServerError(w, r, err)
			return
		}
		if verificationToken != nil {
			go func() {
				// Create a new context as I would like for the email been sent to operate independenty of the request context
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				app.logger.Info("Interaction with notification service begins  ....", *verificationToken)
				n, err := app.notificationService.Send(ctx, &notification.NotificationRequest{
					Email:   user.Email,
					Title:   "Customer Verification",
					Content: fmt.Sprintf("This is your verification token %s", *verificationToken),
				},
				)
				if err != nil {
					app.logger.Warn("Error interacting with the notification service", err)

				} else {

					app.logger.Info("Success interacting with the notification service", n)
				}

			}()
		}
		app.jsonResponse(w, http.StatusCreated, "Customer account created successfully, please check email for a verification link!", user)
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
