package repository

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	OrderID    string
	UserID     string
	BookID     string
	Quantity   int32
	TotalCents int64
	Status     string
	CreatedAt  time.Time
}

type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			order_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			book_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			total_cents INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL
		)
	`)
	return err
}

func (r *PostgresOrderRepo) Create(ctx context.Context, order *Order) error {
	order.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO orders (order_id, user_id, book_id, quantity, total_cents, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		order.OrderID,
		order.UserID,
		order.BookID,
		order.Quantity,
		order.TotalCents,
		order.Status,
		order.CreatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (r *PostgresOrderRepo) UpdateResult(ctx context.Context, orderID string, totalCents int64, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET total_cents = ?, status = ?
		WHERE order_id = ?
	`, totalCents, status, orderID)
	return err
}

func (r *PostgresOrderRepo) ListByUser(ctx context.Context, userID string) ([]Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT order_id, user_id, book_id, quantity, total_cents, status, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		var createdAt string
		if err := rows.Scan(
			&order.OrderID,
			&order.UserID,
			&order.BookID,
			&order.Quantity,
			&order.TotalCents,
			&order.Status,
			&createdAt,
		); err != nil {
			return nil, err
		}

		parsed, err := time.Parse(time.RFC3339Nano, createdAt)
		if err == nil {
			order.CreatedAt = parsed
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
