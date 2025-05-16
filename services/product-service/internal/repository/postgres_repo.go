package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"github.com/lib/pq"
)

type SqlProductRepo struct {
	db *sql.DB
}

func NewPostgresProductRepo(db *sql.DB) *SqlProductRepo {
	return &SqlProductRepo{db}
}

func (s *SqlProductRepo) GetAppProductPolicy(ctx context.Context) (model.AppProductPolicy, error) {
	var policy model.AppProductPolicy
	var formula types.SharingFormula

	query := `
		SELECT 
			p.id,
			p.current_sharing_formula_id,
			p.product_price_to_use,
			sf.id,
			sf.app,
			sf.vendor,
			sf.based_on,
			sf.description
		FROM app_product_policies p
		LEFT JOIN sharing_formulas sf ON p.current_sharing_formula_id = sf.id
		ORDER BY p.id DESC
		LIMIT 1
	`

	err := s.db.QueryRowContext(ctx, query).Scan(
		&policy.Id,
		&policy.CurrentSharingFormulaId,
		&policy.ProductPriceToUse,
		&formula.Id,
		&formula.App,
		&formula.Vendor,
		&formula.BasedOn,
		&formula.Description,
	)

	if err != nil {
		return model.AppProductPolicy{}, err
	}

	policy.CurrentSharingFormula = formula
	return policy, nil
}

func (s *SqlProductRepo) BulkAddProducts(ctx context.Context, payload []ProductInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range payload {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO products (name, description, price, category_label, sub_category_label, tags, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		`, p.Name, p.Description, p.Price, p.CategoryLabel, strings.Join(p.SubCategoryLabel, ","), strings.Join(p.Tags, ","))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SqlProductRepo) DeleteProduct(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (s *SqlProductRepo) UpdateProduct(ctx context.Context, id int, payload ProductInput) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE products
		SET name = $1, description = $2, price = $3, category_label = $4, sub_category_label = $5, tags = $6, updated_at = NOW()
		WHERE id = $7
	`, payload.Name, payload.Description, payload.Price, payload.CategoryLabel, strings.Join(payload.SubCategoryLabel, ","), strings.Join(payload.Tags, ","), id)
	return err
}

func (s *SqlProductRepo) GetProducts(ctx context.Context, pagination *utils.PaginationPayload) ([]model.Product, int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, description, price, tags, created_at, updated_at
		FROM products
		LIMIT $2 OFFSET $3
	`, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		var tags string
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price.Amount, &tags, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		p.Tags = strings.Split(tags, ",")
		products = append(products, p)
	}

	// count total
	var total int
	s.db.QueryRowContext(ctx, `SELECT count(*) FROM products`).Scan(&total)
	return products, total, nil
}

func (s *SqlProductRepo) UpdateProductInventory(ctx context.Context, id int, storeId int, productId int, quantity int, metaData *map[string]string) error {
	var metaDataJSON []byte
	var err error
	if metaData != nil {
		metaDataJSON, err = json.Marshal(metaData)
		if err != nil {
			return err
		}
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO inventory (id, store_id, product_id, quantity, meta_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (id, store_id, product_id)
		DO UPDATE SET
			quantity = EXCLUDED.quantity,
			meta_data = EXCLUDED.meta_data,
			updated_at = NOW()
	`, id, storeId, productId, quantity, metaDataJSON)

	return err
}

func (s *SqlProductRepo) BulkAddCategories(ctx context.Context, payload []CategoryInput) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, c := range payload {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO categories (name, description, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
		`, c.Name, c.Description)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SqlProductRepo) DeleteCategory(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}

func (s *SqlProductRepo) UpdateCategory(ctx context.Context, id int, payload CategoryInput) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`, payload.Name, payload.Description, id)
	return err
}

func (s *SqlProductRepo) GetCategories(ctx context.Context, pagination *utils.PaginationPayload) ([]model.Category, int, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		LIMIT $2 OFFSET $3
	`, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	// count total
	var total int
	s.db.QueryRowContext(ctx, `SELECT count(*) FROM categories`).Scan(&total)
	return categories, total, nil
}

func (s *SqlProductRepo) CreateSharingFormula(ctx context.Context, id int, basedOn types.SharingFormulaBasedOn, appPercent int, vendorPercent int, description string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sharing_formulas (id, based_on, app, vendor, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`, id, basedOn, appPercent, vendorPercent, description)
	return err
}

func (s *SqlProductRepo) SaveAppProductPolicy(ctx context.Context, sharingFormulaId int, priceToUse types.DominantPriceType) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE app_product_policy
		SET sharing_formula_id = $1, product_price_to_use = $2, updated_at = NOW()
		WHERE id = 1
	`)
	return err
}

func (s *SqlProductRepo) CreateDiscount(ctx context.Context, payload model.Discount) error {
	applicableToJson, err := json.Marshal(payload.ApplicableTo)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO discounts (id, value, type, effective_at, expires_at, paid_by, applicable_to, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`, payload.Id, payload.Value, payload.ValueType, payload.EffectiveAt, payload.ExpiresAt, payload.PaidBy, applicableToJson, payload.Name, payload.Description)
	return err
}

func (s *SqlProductRepo) UpdateDiscountApplicability(ctx context.Context, id int, payload types.DiscountApplicability) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE discounts SET applicable_to = $1, updated_at = NOW() WHERE id = $2
	`, data, id)
	return err
}

func (s *SqlProductRepo) UpdateDiscountExpiryDate(ctx context.Context, id int, expiryDate time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE discounts SET expires_at = $1, updated_at = NOW() WHERE id = $2
	`, expiryDate, id)
	return err
}

func (s *SqlProductRepo) GetDiscounts(ctx context.Context, pagination *utils.PaginationPayload, filter *DiscountFilter) (result []model.Discount, total int, err error) {
	args := []interface{}{}
	whereClauses := []string{}
	joinClause := ""

	argIndex := 1

	// Handle ExpiresAt filter
	if filter != nil && filter.ExpiresAt != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("d.expires_at >= $%d", argIndex))
		args = append(args, *filter.ExpiresAt)
		argIndex++
	}

	// Handle Applicability filter (assumes normalized discount_applicabilities table)
	if filter != nil && (len(filter.Applicability.ProductIds) > 0 ||
		len(filter.Applicability.StoreProductIds) > 0 ||
		len(filter.Applicability.StoreProductInventoryIds) > 0) {

		joinClause = "JOIN discount_applicabilities da ON da.discount_id = d.id"

		applicabilityConditions := []string{}

		if len(filter.Applicability.ProductIds) > 0 {
			applicabilityConditions = append(applicabilityConditions, fmt.Sprintf("da.product_id = ANY($%d)", argIndex))
			args = append(args, pq.Array(filter.Applicability.ProductIds))
			argIndex++
		}

		if len(filter.Applicability.StoreProductIds) > 0 {
			applicabilityConditions = append(applicabilityConditions, fmt.Sprintf("da.store_product_id = ANY($%d)", argIndex))
			args = append(args, pq.Array(filter.Applicability.StoreProductIds))
			argIndex++
		}

		if len(filter.Applicability.StoreProductInventoryIds) > 0 {
			applicabilityConditions = append(applicabilityConditions, fmt.Sprintf("da.store_product_inventory_id = ANY($%d)", argIndex))
			args = append(args, pq.Array(filter.Applicability.StoreProductInventoryIds))
			argIndex++
		}

		whereClauses = append(whereClauses, "("+strings.Join(applicabilityConditions, " OR ")+")")
	}

	// Pagination
	args = append(args, pagination.Limit, pagination.Offset)

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.name, d.description, d.value, d.value_type, d.effective_at, d.expires_at, d.paid_by, d.created_at, d.updated_at
		FROM 
			discounts d
			%s
		%s
		ORDER BY d.created_at DESC
		LIMIT $%d OFFSET $%d
	`, joinClause, whereSQL, argIndex, argIndex+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var discounts []model.Discount
	for rows.Next() {
		var d model.Discount
		err := rows.Scan(
			&d.Id,
			&d.Name,
			&d.Description,
			&d.Value,
			&d.ValueType,
			&d.EffectiveAt,
			&d.ExpiresAt,
			&d.PaidBy,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		discounts = append(discounts, d)
	}

	// Total count (filtered)
	countQuery := fmt.Sprintf(`SELECT COUNT(DISTINCT d.id) FROM discounts d %s %s`, joinClause, whereSQL)
	err = s.db.QueryRowContext(ctx, countQuery, args[:argIndex-1]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}
