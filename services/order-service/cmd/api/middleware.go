package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/kaasikodes/shop-ease/shared/proto/auth"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) isCustomerActiveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := app.trace.Start(ctx, "Customer middleware")
		defer span.End()
		userId, ok := getUserIdFromContext(ctx)
		if !ok {
			err := errors.New("unable to retrieve userId")
			app.logger.WithContext(ctx).Error("Unable to retrieve userId", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			app.badRequestResponse(w, r, err)
			return
		}

		// get user from auth service
		user, err := app.clients.auth.GetUserById(ctx, &auth.GetUserByIdRequest{UserId: int32(userId)})
		if err != nil {
			app.logger.WithContext(ctx).Error("err", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("unable to retrieve client: %w", err))
			return
		}
		userRoles := user.User.Roles
		isCustomerActive := false
		for _, role := range userRoles {
			if role.IsActive && role.Name == "customer" {
				isCustomerActive = true
				return
			}
		}
		app.logger.WithContext(ctx).Info("user info for customer", user)

		if !isCustomerActive {
			app.unauthorizedErrorResponse(w, r, errors.New("customer account is deactivated"))
			return

		}

		// Step 5: Call next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.trace.Start(r.Context(), "Auth middleware")
		defer span.End()

		// Step 1: Verify token
		claims, err := app.jwt.ExtractAndVerifyToken(r)
		if err != nil {
			app.logger.WithContext(ctx).Error("JWT token error", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("please provide a valid token: %w", err))
			return
		}

		// Step 2: Retrieve user from DB
		userId, err := strconv.Atoi(claims.UserID)
		if err != nil {
			app.logger.WithContext(ctx).Error("Invalid user ID in token", err)
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid user ID in token"))
			return
		}

		// Step 4: Add user and claims to context
		ctx = context.WithValue(ctx, ContextKeyUser{}, userId)

		// Step 5: Call next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		normalizedPath := normalizePath(r.URL.Path)
		// path := r.URL.Path
		path := normalizedPath //normalize paths (e.g., replace /v1/users/123 with /v1/users/:id) to reduce high cardinality in metrics.
		method := r.Method

		start := time.Now()

		// In-flight request increment
		app.metrics.inFlightRequests.Inc()
		defer app.metrics.inFlightRequests.Dec()

		// Capture request size (from Content-Length)
		reqSize := float64(0)
		if r.ContentLength > 0 {
			reqSize = float64(r.ContentLength)
		}
		app.metrics.requestSize.WithLabelValues(path, method).Observe(reqSize)

		// Wrap the ResponseWriter to capture status code and response size
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: 200}

		// Call next handler
		next.ServeHTTP(rw, r)

		// Observe duration
		duration := time.Since(start).Seconds()
		app.metrics.requestDuration.WithLabelValues(path, method).Observe(duration)

		// Count response size
		app.metrics.responseSize.WithLabelValues(path, method).Observe(float64(rw.size))

		// Count requests
		app.metrics.requestCount.WithLabelValues(path, method).Inc()

		// Count errors
		if rw.statusCode >= 400 {
			app.metrics.errorCount.WithLabelValues(path, method, strconv.Itoa(rw.statusCode)).Inc()
		}
	})
}
func normalizePath(path string) string {
	// You can add more normalization patterns as needed
	userIdPattern := regexp.MustCompile(`/v1/users/\d+`)
	orderIdPattern := regexp.MustCompile(`/v1/orders/\d+`)
	uuidPattern := regexp.MustCompile(`/[0-9a-fA-F\-]{36}`)

	// Apply them in order
	path = userIdPattern.ReplaceAllString(path, "/v1/users/:id")
	path = orderIdPattern.ReplaceAllString(path, "/v1/orders/:id")
	path = uuidPattern.ReplaceAllString(path, "/:uuid")

	return path
}

// Custom response writer to capture response size and status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
