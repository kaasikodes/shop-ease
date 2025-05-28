package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
)

type SqlPaymentRepo struct {
	db *sql.DB
}

func NewSqlPaymentRepo(db *sql.DB) *SqlPaymentRepo {
	return &SqlPaymentRepo{db}

}

func (p *SqlPaymentRepo) CreateTransaction(tx model.Transaction) (*model.Transaction, error) {
	const query = `
		INSERT INTO transactions (provider, transaction_id, meta_data, entity_id, amount, entity_payment_type, status, paid_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	metaDataJson, err := json.Marshal(tx.MetaData)
	if err != nil {
		return nil, err
	}

	result, err := p.db.Exec(query,
		tx.Provider,
		tx.TransactionId,
		string(metaDataJson),
		tx.EntityId,
		tx.Amount,
		tx.EntityPaymentType,
		tx.Status,
		tx.PaidAt,
	)
	if err != nil {
		return nil, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	tx.ID = int(insertedID)
	return &tx, nil
}

func (p *SqlPaymentRepo) GetTransactionById(id int) (*model.Transaction, error) {
	const query = `
		SELECT id, provider, transaction_id, meta_data, entity_id, amount, entity_payment_type, status, paid_at, created_at, updated_at
		FROM transactions WHERE id = ?
	`

	var tx model.Transaction
	var metaDataStr string
	err := p.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.Provider,
		&tx.TransactionId,
		&metaDataStr,
		&tx.EntityId,
		&tx.Amount,
		&tx.EntityPaymentType,
		&tx.Status,
		&tx.PaidAt,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(metaDataStr), &tx.MetaData)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (p *SqlPaymentRepo) UpdateTransaction(id int, payload model.Transaction) (*model.Transaction, error) {
	const query = `
		UPDATE transactions
		SET provider = ?, transaction_id = ?, meta_data = ?, entity_id = ?, amount = ?, 
		    entity_payment_type = ?, status = ?, paid_at = ?, updated_at = NOW()
		WHERE id = ?
	`

	metaDataJson, err := json.Marshal(payload.MetaData)
	if err != nil {
		return nil, err
	}

	_, err = p.db.Exec(query,
		payload.Provider,
		payload.TransactionId,
		string(metaDataJson),
		payload.EntityId,
		payload.Amount,
		payload.EntityPaymentType,
		payload.Status,
		payload.PaidAt,
		id,
	)
	if err != nil {
		return nil, err
	}

	payload.ID = id
	return &payload, nil
}

func (p *SqlPaymentRepo) GetTransactions(pagination *types.PaginationPayload, filter *model.TransactionFilter) ([]model.Transaction, int, error) {
	var filters []string
	var args []interface{}

	if filter != nil {
		if filter.Provider != "" {
			filters = append(filters, "provider = ?")
			args = append(args, filter.Provider)
		}
		if filter.Status != "" {
			filters = append(filters, "status = ?")
			args = append(args, filter.Status)
		}
		if filter.EntityPaymentType != "" {
			filters = append(filters, "entity_payment_type = ?")
			args = append(args, filter.EntityPaymentType)
		}
		if filter.PaidAt != nil {
			filters = append(filters, "DATE(paid_at) = DATE(?)")
			args = append(args, filter.PaidAt)
		}
		if filter.Amount > 0 {
			filters = append(filters, "amount = ?")
			args = append(args, filter.Amount)
		}
	}

	whereClause := ""
	if len(filters) > 0 {
		whereClause = "WHERE " + strings.Join(filters, " AND ")
	}

	limit := pagination.Limit
	offset := (pagination.Offset - 1) * limit

	query := fmt.Sprintf(`
		SELECT id, provider, transaction_id, meta_data, entity_id, amount, entity_payment_type, status, paid_at, created_at, updated_at
		FROM transactions
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	argsWithPagination := append(args, limit, offset)

	rows, err := p.db.Query(query, argsWithPagination...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		var metaDataStr string
		err := rows.Scan(
			&tx.ID,
			&tx.Provider,
			&tx.TransactionId,
			&metaDataStr,
			&tx.EntityId,
			&tx.Amount,
			&tx.EntityPaymentType,
			&tx.Status,
			&tx.PaidAt,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		json.Unmarshal([]byte(metaDataStr), &tx.MetaData)
		results = append(results, tx)
	}

	// Total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions %s`, whereClause)
	var total int
	err = p.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
