package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func openOrderTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestCreateRejectsInvalidOrder(t *testing.T) {
	repo := NewPostgresOrderRepo(openOrderTestDB(t))
	err := repo.Create(context.Background(), &Order{
		OrderID: "", UserID: "u-100", BookID: "b-100", Quantity: 1,
	})
	if err == nil {
		t.Fatal("expected invalid order to be rejected")
	}
}

func TestUpdateResultRejectsUnknownStatus(t *testing.T) {
	repo := NewPostgresOrderRepo(openOrderTestDB(t))
	err := repo.UpdateResult(context.Background(), "ord-100", 2999, "unknown")
	if err == nil {
		t.Fatal("expected unsupported status to be rejected")
	}
}

func TestListByUserRejectsMissingUser(t *testing.T) {
	repo := NewPostgresOrderRepo(openOrderTestDB(t))
	if _, err := repo.ListByUser(context.Background(), " "); err == nil {
		t.Fatal("expected missing user to be rejected")
	}
}
