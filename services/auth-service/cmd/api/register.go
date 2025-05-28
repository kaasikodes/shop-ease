package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"github.com/kaasikodes/shop-ease/services/payment-service/pkg/model"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"github.com/kaasikodes/shop-ease/shared/proto/payment"
	"github.com/kaasikodes/shop-ease/shared/proto/subscription"
	"github.com/kaasikodes/shop-ease/shared/proto/vendor_service"
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
	SubscriptionPlanId int64
	Store              struct {
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

	switch payload.RoleId {
	case store.CustomerID:
		app.logger.Info("New Customer Registeration initiated ...")
		user, verificationToken, err := app.registerCustomer(ctx, payload, false)
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
		paymentUrl, err := app.registerVendor(ctx, payload)
		if err != nil {
			app.logger.WithContext(registerTraceCtx).Error("Error registering vendor", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			if err == store.ErrDuplicateEmail {
				app.badRequestResponse(w, r, err)
				return
			}
			app.internalServerError(w, r, err)
			return
		}
		app.jsonResponse(w, http.StatusCreated, "Please use the payment link to pay for your vendor subscription!", map[string]string{"paymentUrl": paymentUrl})
		return
	default:
		app.badRequestResponse(w, r, errors.New("please select a valid role id"))
		return

	}

}
func (app *application) registerCustomer(ctx context.Context, payload RegisterUserPayload, isVerifiedByOauth bool) (*store.User, *string, error) {
	_, span := app.trace.Start(ctx, "register-customer")
	defer span.End()
	var plainToken string
	// create user acc
	user := &store.User{
		Name:  payload.Name,
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		return nil, nil, err
	}
	if isVerifiedByOauth {

		// Start a new transaction
		tx, err := app.store.BeginTx(ctx)
		if err != nil {
			span.RecordError(err)
			return nil, nil, err
		}

		err = app.store.Users().Create(ctx, tx, user, &store.UserRole{
			ID: store.CustomerID,
		})
		if err != nil {
			span.RecordError(err)
			tx.Rollback()
			return nil, nil, err
		}
		err = app.store.Users().Verify(ctx, tx, user)
		if err != nil {
			span.RecordError(err)
			tx.Rollback()
			return nil, nil, err
		}
		// Commit the transaction
		if err := tx.Commit(); err != nil {
			span.RecordError(err)
			return nil, nil, err
		}

		return user, &plainToken, nil

	}

	plainToken = uuid.New().String() //TODO: hash and save token, not just the plain token
	err := app.store.Users().CreateWithVerificationToken(ctx, user, plainToken, ExpiresAtVerificationToken)

	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}

	// communicate with the notification service via grpc to send a verification email
	return user, &plainToken, nil

}
func (app *application) registerVendor(ctx context.Context, payload RegisterUserPayload) (paymentUrl string, err error) {
	// TODO: Not enough use cases accounted for, for verification on payment should verify user if not verified
	_, span := app.trace.Start(ctx, "register-vendor")
	defer span.End()
	user := &store.User{
		Name:  payload.Name,
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		return "", err
	}

	// Start a new transaction
	tx, err := app.store.BeginTx(ctx)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	err = app.store.Users().Create(ctx, tx, user, &store.UserRole{
		ID: store.CustomerID,
	})
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return "", err
	}
	err = app.store.Users().Verify(ctx, tx, user)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return "", err
	}

	// create vendor
	vendor, err := app.clients.vendor.CreateVendor(ctx, &vendor_service.CreateVendorRequest{
		Email:  payload.Email,
		Name:   payload.Name,
		Phone:  &payload.Vendor.Store.Contact.phone,
		UserId: int64(user.ID),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		tx.Rollback()

		return "", nil
	}
	// create subscription plan for vendor
	subscription, err := app.clients.subscription.CreateVendorSubscription(ctx, &subscription.CreateVendorSubscriptionRequest{
		VendorId: vendor.Id,
		PlanId:   payload.Vendor.SubscriptionPlanId,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		tx.Rollback()

		return "", nil
	}
	// then create payment for subscription
	paymentData, err := app.clients.payment.CreateTransaction(ctx, &payment.CreateTransactionRequest{
		EntityId:          subscription.Id,
		Provider:          string(model.PaymentProviderPaystack),
		EntityPaymentType: string(model.EntityPaymentTypeVendorSubscriptionPayment),
		Amount:            subscription.Plan.Amount,
		MetaData: map[string]string{
			"reason":      "first time payment for vendor subscription; new vendor registration",
			"amount_paid": strconv.Itoa(int(subscription.Plan.Amount)),
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		tx.Rollback()

		return "", nil
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		span.RecordError(err)
		return "", err
	}
	// return the payment url to user
	return paymentData.PaymentUrl, nil

	// the payload should account for the subscription plan selected as a vendor,
	// a vendor acccount should be created (vendor service will have a middleware that will always check with subscription service wether or not the the vendor_sub_plan_is_active)
	// Talk to subscription service to update the vendors on this plan
	// subscription service should contain vendors that belong to a subscription and when their subscription expires, subscription should have instances that are tied to vendors in the event that a subscription plan is modified, - emit subscription created
	// now the payment service will contain the webhook and when informed of payment will inform the sunscription service that payment is complete
	// on receing that subscrption service updates the vendor subscription
	// now each the vendor will have a middleware that checks wether the vendor is active,
	// now on subscription.vendor_plan_renewed the vendor service is notified to activate vendor account, and on subscription.vendor_plan_expired the vendor service is deactivate

	// FLow
	// a if user already exists that is fine, as will be uspsert expect for the subscription service where the rule is a vendor can have only one active plan(enforce on db level[not possible, so on code level, and add middleware to perform check])
	// Dont send a verification until payment is made,and on payment is made regardles of the invalidity of mail or whatever verify account
	// auth -> vendor service (create inactive vendor, with the user_id specified) -> subscription (create an instance plan for vendor -snapshot  of plan at the moment, amount_to_be paid, plan_id, vendor_id, isActive(yes or no), paidAt, isPaid, CU, expiredAt, ). -> payment (receives the sub_plan_instance_id, and the amount, and creates a payment link): Auth then Returns a message to client with link and message of user to complete payment with link provided. Payment service has webhook with the gateway concerned and on confirmation will emit payment.vendor_plan_paid and subscriptio will listen and update the plan_instance_id, notification service will listen as well and notify the user that the vendor account is activated

}

// customer path - enter email, name, and password. Gets a verification mail via the notification service. The user clicks on the link on the verification mail(will contain a token) - token is valid, and then verified, and then account is verified & customer role is active or the token is invalid and an error message is sent to the user. Ensure validation shows proper error messages and path and all. All these events will be logged and how can this be used for data analysis, what role as well can prometheus and grafana play in this.
// vendor path
// a bit more info, and will have to interact with the subscription service first
