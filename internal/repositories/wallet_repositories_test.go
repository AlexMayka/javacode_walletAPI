package repositories_test

import (
	"JavaCode/internal/repositories"
	"JavaCode/utils"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"testing"
	"time"
)

func TestGetWalletByUUID(t *testing.T) {
	t.Run("Test 1: Wallet found", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		now := time.Now()

		mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1").
			WithArgs(walletID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
					AddRow(walletID, 1000, now, now),
			)

		result, err := repositories.GetWalletByUUID(db, walletID)
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}

		if result.Id != walletID || result.Balance != 1000 {
			t.Errorf("unexpected wallet data: %+v", result)
		}
	})

	t.Run("Test 2: Wallet not found", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "not-found"

		mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1").
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)

		_, err := repositories.GetWalletByUUID(db, walletID)
		if !errors.Is(err, utils.ErrWalletNotFound) {
			t.Errorf("expected ErrWalletNotFound, got: %v", err)
		}
	})

	t.Run("Test 3: Other DB error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"

		mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1").
			WithArgs(walletID).
			WillReturnError(sql.ErrConnDone)

		_, err := repositories.GetWalletByUUID(db, walletID)
		if !errors.Is(err, sql.ErrConnDone) {
			t.Errorf("expected sql.ErrConnDone, got: %v", err)
		}
	})
}

func TestGetWalletForUpdate(t *testing.T) {
	t.Run("Test 1: Wallet found and locked", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		now := time.Now()

		mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
				AddRow(walletID, 1500, now, now))

		result, err := repositories.GetWalletForUpdate(db, walletID)
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
		if result.Id != walletID || result.Balance != 1500 {
			t.Errorf("unexpected result: %+v", result)
		}
	})

	t.Run("Test 2: Wallet not found", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "not-found"

		mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1 FOR UPDATE").
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)

		_, err := repositories.GetWalletForUpdate(db, walletID)
		if !errors.Is(err, utils.ErrWalletNotFound) {
			t.Errorf("expected ErrWalletNotFound, got: %v", err)
		}
	})
}

func TestChainBalance(t *testing.T) {
	t.Run("Test 1: Successful balance update", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		delta := 500

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

		err := repositories.ChainBalance(db, walletID, delta)
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
	})

	t.Run("Test 2: Wallet not found (0 rows affected)", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "not-found"
		delta := 100

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // no rows updated

		err := repositories.ChainBalance(db, walletID, delta)
		if !errors.Is(err, utils.ErrWalletNotFound) {
			t.Errorf("expected ErrWalletNotFound, got: %v", err)
		}
	})

	t.Run("Test 3: Constraint violation (balance < 0)", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		delta := -99999

		pqErr := &pq.Error{Constraint: "wallets_balance_check"}

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnError(pqErr)

		err := repositories.ChainBalance(db, walletID, delta)
		if !errors.Is(err, utils.ErrNegativeBalance) {
			t.Errorf("expected ErrNegativeBalance, got: %v", err)
		}
	})

	t.Run("Test 4: Generic SQL error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		delta := 100

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnError(sql.ErrConnDone)

		err := repositories.ChainBalance(db, walletID, delta)
		if !errors.Is(err, sql.ErrConnDone) {
			t.Errorf("expected sql.ErrConnDone, got: %v", err)
		}
	})

	t.Run("Test 5: Error during RowsAffected()", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		walletID := "abc-123"
		delta := 100

		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		err := repositories.ChainBalance(db, walletID, delta)
		if err == nil || err.Error() != "rows affected error" {
			t.Errorf("expected rows affected error, got: %v", err)
		}
	})
}
