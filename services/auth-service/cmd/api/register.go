package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaasikodes/shop-ease/internal/store"

	"github.com/google/uuid"
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
	ctx := r.Context()
	switch payload.RoleId {
	case store.CustomerID:
		user, err := app.registerCustomer(ctx, payload)
		if err != nil {
			if err == store.ErrDuplicateEmail {
				app.badRequestResponse(w, r, err)
				return

			}
			app.internalServerError(w, r, err)
			return
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
func (app *application) registerCustomer(ctx context.Context, payload RegisterUserPayload) (*store.User, error) {
	// create user acc
	user := &store.User{
		Name:  payload.Name,
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		return nil, err
	}
	plainToken := uuid.New().String() //TODO: hash and save token, not just the plain token
	err := app.store.Users().CreateWithVerificationToken(ctx, user, plainToken, ExpiresAtVerificationToken)

	if err != nil {
		return nil, err
	}

	// communicate with the notification service via grpc to send a verification email
	return user, nil

}
func registerVendor() {

}

// customer path - enter email, name, and password. Gets a verification mail via the notification service. The user clicks on the link on the verification mail(will contain a token) - token is valid, and then verified, and then account is verified & customer role is active or the token is invalid and an error message is sent to the user. Ensure validation shows proper error messages and path and all. All these events will be logged and how can this be used for data analysis, what role as well can prometheus and grafana play in this.
// vendor path
// a bit more info, and will have to interact with the subscription service first
