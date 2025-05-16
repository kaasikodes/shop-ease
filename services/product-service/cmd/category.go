package main

import (
	"errors"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Categories")

	defer span.End()

	app.logger.WithContext(initialTraceCtx).Info("getting categories")

	pagination := utils.GetPaginationFromQuery(r)

	categories, total, err := app.store.GetCategories(initialTraceCtx, &utils.PaginationPayload{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	span.SetAttributes(
		attribute.Int("pagination.limit", pagination.Limit),
		attribute.Int("pagination.offset", pagination.Offset),
	)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting categories", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting categories")

	var result = make([]any, len(categories))
	for i, category := range categories {
		result[i] = category

	}

	app.jsonResponse(w, http.StatusOK, "Categories retrieved successfully!", createPaginatedResponse(result, total))
	return

}

func (app *application) bulkAddCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Bulk Add Categories")
	defer span.End()

	var input []repository.CategoryInput
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(input) == 0 {
		app.badRequestResponse(w, r, errors.New("category list cannot be empty"))
		return
	}

	err = app.store.BulkAddCategories(ctx, input)
	if err != nil {
		app.logger.WithContext(ctx).Error("Error adding categories", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, "Categories added successfully", nil)
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Delete Category")
	defer span.End()

	categoryID, err := app.readIntParam(r, "categoryId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	err = app.store.DeleteCategory(ctx, categoryID)
	if err != nil {
		app.logger.WithContext(ctx).Error("Error deleting category", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Category deleted successfully", nil)
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Update Category")
	defer span.End()

	categoryID, err := app.readIntParam(r, "categoryId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	var input repository.CategoryInput
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.UpdateCategory(ctx, categoryID, input)
	if err != nil {
		app.logger.WithContext(ctx).Error("Error updating category", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Category updated successfully", nil)
}
