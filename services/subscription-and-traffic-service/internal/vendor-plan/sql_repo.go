package vendorplan

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type SqlVendorRepo struct {
	db *sql.DB
}

func NewSqlVendorRepo(db *sql.DB) *SqlVendorRepo {
	return &SqlVendorRepo{db}

}

func (r *SqlVendorRepo) GetVendorUserInteractionRecords(pagination *types.PaginationPayload, filter *VendorUserInteractionFilter) (result []VendorUserInteraction, total int, err error) {
	var conditions []string
	var args []interface{}

	if filter != nil {
		if filter.VendorId != 0 {
			conditions = append(conditions, "vendor_id = ?")
			args = append(args, filter.VendorId)
		}
		if filter.UserId != 0 {
			conditions = append(conditions, "user_id = ?")
			args = append(args, filter.UserId)
		}
		if filter.Type != "" {
			conditions = append(conditions, "type = ?")
			args = append(args, filter.Type)
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vendor_user_interactions %s", whereClause)
	err = r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Select records with pagination
	query := fmt.Sprintf(`
		SELECT id, vendor_id, user_id, type, created_at, updated_at 
		FROM vendor_user_interactions %s 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec VendorUserInteraction
		err = rows.Scan(&rec.ID, &rec.VendorId, &rec.UserId, &rec.Type, &rec.CreatedAt, &rec.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, rec)
	}

	return result, total, nil
}

func (r *SqlVendorRepo) GetVendorPlans(pagination *types.PaginationPayload, filter *VendorPlanFilter) (result []VendorPlan, total int, err error) {
	var conditions []string
	var args []interface{}

	if filter != nil {
		if filter.IsActive != nil {
			conditions = append(conditions, "is_active = ?")
			args = append(args, *filter.IsActive)
		}
		if filter.Name != "" {
			conditions = append(conditions, "name LIKE ?")
			args = append(args, "%"+filter.Name+"%")
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vendor_plans %s", whereClause)
	err = r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, name, content, price, user_interactions_allowed, duration_in_secs, is_active, created_at, updated_at 
		FROM vendor_plans %s 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var plan VendorPlan
		err = rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Content,
			&plan.Price,
			&plan.UserInteractionsAllowed,
			&plan.DurationInSecs,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, plan)
	}

	return result, total, nil
}

func (r *SqlVendorRepo) CreateVendorPlan(payload VendorPlanPayload) (*VendorPlan, error) {
	query := `
		INSERT INTO vendor_plans 
		(name, content, price, user_interactions_allowed, duration_in_secs, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, TRUE, NOW(), NOW())
	`
	res, err := r.db.Exec(query, payload.Name, payload.Content, int(payload.Price), payload.UserInteractionsAllowed, int64(payload.DurationInSecs.Seconds()))
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	plan := &VendorPlan{}
	err = r.db.QueryRow("SELECT id, name, content, price, user_interactions_allowed, duration_in_secs, is_active, created_at, updated_at FROM vendor_plans WHERE id = ?", id).
		Scan(
			&plan.ID,
			&plan.Name,
			&plan.Content,
			&plan.Price,
			&plan.UserInteractionsAllowed,
			&plan.DurationInSecs,
			&plan.IsActive,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (r *SqlVendorRepo) BulkActivateOrDeactivateVendorPlan(planIds []int, isActive bool) error {
	if len(planIds) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(planIds))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("UPDATE vendor_plans SET is_active = ?, updated_at = NOW() WHERE id IN (%s)", placeholders)

	args := []interface{}{isActive}
	for _, id := range planIds {
		args = append(args, id)
	}

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SqlVendorRepo) CreateVendorPlanSubscription(planId int, vendorId int) (*VendorSubsription, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(30*24) * time.Hour) // default 30 days, adjust if needed

	query := `
		INSERT INTO vendor_subscriptions 
		(plan_id, vendor_id, has_paid, limit_exceeded_at, paid_at, began_at, expires_at, created_at, updated_at)
		VALUES (?, ?, FALSE, NULL, ?, ?, ?, NOW(), NOW())
	`

	res, err := r.db.Exec(query, planId, vendorId, now, now, expiresAt)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	sub := &VendorSubsription{}
	err = r.db.QueryRow(`
		SELECT id, plan_id, vendor_id, has_paid, limit_exceeded_at, paid_at, began_at, expires_at, created_at, updated_at 
		FROM vendor_subscriptions WHERE id = ?`, id).
		Scan(&sub.ID, &sub.PlanId, &sub.VendorId, &sub.HasPaid, &sub.LimitExceededAt, &sub.PaidAt, &sub.BeganAt, &sub.ExpiresAt, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *SqlVendorRepo) UpdateVendorPlanSubscription(subscriptionId int, payload VendorSubsription) (*VendorSubsription, error) {
	// Prepare the SQL update statement
	query := `
		UPDATE vendor_subscriptions
		SET plan_id = ?, vendor_id = ?, has_paid = ?, limit_exceeded_at = ?, paid_at = ?, began_at = ?, expires_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, plan_id, vendor_id, has_paid, limit_exceeded_at, paid_at, began_at, expires_at, created_at, updated_at
	`

	row := r.db.QueryRow(
		query,
		payload.PlanId,
		payload.VendorId,
		payload.HasPaid,
		payload.LimitExceededAt,
		payload.PaidAt,
		payload.BeganAt,
		payload.ExpiresAt,
		subscriptionId,
	)

	var updatedSub VendorSubsription
	err := row.Scan(
		&updatedSub.ID,
		&updatedSub.PlanId,
		&updatedSub.VendorId,
		&updatedSub.HasPaid,
		&updatedSub.LimitExceededAt,
		&updatedSub.PaidAt,
		&updatedSub.BeganAt,
		&updatedSub.ExpiresAt,
		&updatedSub.CreatedAt,
		&updatedSub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedSub, nil
}

func (r *SqlVendorRepo) CreateVendorUserInteractionRecord(userId int, interactionType VendorUserInteractionType) (*VendorUserInteraction, error) {
	query := `
		INSERT INTO vendor_user_interactions (user_id, type, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, vendor_id, user_id, type, created_at, updated_at
	`

	// Since VendorId is missing from function args, I'm assuming it might be inferred from userId or ignored
	// But your struct requires VendorId. So let's assume VendorId = 0 for now. You can modify this as needed.

	var interaction VendorUserInteraction
	err := r.db.QueryRow(query, userId, string(interactionType)).Scan(
		&interaction.ID,
		&interaction.VendorId,
		&interaction.UserId,
		&interaction.Type,
		&interaction.CreatedAt,
		&interaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &interaction, nil
}
func (r *SqlVendorRepo) GetActiveSubscriptionsForVendor(vendorId int64) ([]*VendorSubsription, error) {
	const query = `
		SELECT 
			id, plan_id, vendor_id, has_paid, limit_exceeded_at, paid_at, began_at, expires_at,
			created_at, updated_at
		FROM vendor_subscriptions
		WHERE vendor_id = ? 
		  AND expires_at > NOW() 
		  AND (limit_exceeded_at IS NULL OR limit_exceeded_at > NOW())
	`

	rows, err := r.db.Query(query, vendorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*VendorSubsription

	for rows.Next() {
		var sub VendorSubsription
		var limitExceededAt, paidAt sql.NullTime

		err = rows.Scan(
			&sub.ID,
			&sub.PlanId,
			&sub.VendorId,
			&sub.HasPaid,
			&limitExceededAt,
			&paidAt,
			&sub.BeganAt,
			&sub.ExpiresAt,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if limitExceededAt.Valid {
			sub.LimitExceededAt = limitExceededAt.Time
		}
		if paidAt.Valid {
			sub.PaidAt = paidAt.Time
		}

		subs = append(subs, &sub)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *SqlVendorRepo) GetVendorSubscriptionID(subscriptionId int64) (*VendorSubsription, error) {
	const query = `
		SELECT 
			id, plan_id, vendor_id, has_paid, limit_exceeded_at, paid_at, began_at, expires_at,
			created_at, updated_at
		FROM vendor_subscriptions
		WHERE id = ?
	`

	var sub VendorSubsription
	var limitExceededAt, paidAt sql.NullTime

	err := r.db.QueryRow(query, subscriptionId).Scan(
		&sub.ID,
		&sub.PlanId,
		&sub.VendorId,
		&sub.HasPaid,
		&limitExceededAt,
		&paidAt,
		&sub.BeganAt,
		&sub.ExpiresAt,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // subscription not found
		}
		return nil, err
	}

	if limitExceededAt.Valid {
		sub.LimitExceededAt = limitExceededAt.Time
	}
	if paidAt.Valid {
		sub.PaidAt = paidAt.Time
	}

	return &sub, nil
}

func (r *SqlVendorRepo) MarkSubscriptionPaid(subscriptionId int64) error {
	const query = `
		UPDATE vendor_subscriptions
		SET has_paid = TRUE,
		    paid_at = NOW(),
		    updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.Exec(query, subscriptionId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found with id %d", subscriptionId)
	}

	return nil
}
