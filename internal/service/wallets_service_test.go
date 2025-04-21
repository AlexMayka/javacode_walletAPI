package service_test

import (
	"JavaCode/internal/service"
	"JavaCode/utils"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
)

func expectTxWithBalance(mock sqlmock.Sqlmock, walletID string, balance int, delta int, execErr error) {
	mock.ExpectBegin()

	mock.ExpectQuery("SELECT id, balance.*FOR UPDATE").
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow(walletID, balance, time.Now(), time.Now()))

	if execErr != nil {
		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnError(execErr)
		mock.ExpectRollback()
	} else {
		mock.ExpectExec("UPDATE wallets SET balance = balance \\+ \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
			WithArgs(delta, walletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}
}

func TestGetWalletsService(t *testing.T) {
	t.Run("Test 1: Not find wallet", func(t *testing.T) {
		test := "f4c863ec-0300-495d-852d-c115e197390b"

		db, mock, _ := sqlmock.New()
		defer db.Close()

		q := "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1"
		mock.ExpectQuery(q).WithArgs(test).WillReturnError(sql.ErrNoRows)

		_, err := service.GetWalletsService(db, test)

		if !errors.Is(err, utils.ErrWalletNotFound) {
			t.Errorf("TestGetWalletsService: got %v, want %v", err, utils.ErrWalletNotFound)
		}
	})

	t.Run("Test 2: Error DataBase", func(t *testing.T) {
		test := "f4c863ec-0300-495d-852d-c115e197390b"

		db, mock, _ := sqlmock.New()
		defer db.Close()

		q := "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1"
		mock.ExpectQuery(q).WithArgs(test).WillReturnError(sql.ErrConnDone)

		_, err := service.GetWalletsService(db, test)

		t.Log(err.Error())

		if !errors.Is(err, utils.ErrDatabase) {
			t.Errorf("TestGetWalletsService: got %v, want %v", err, utils.ErrDatabase)
		}
	})

	t.Run("Test 3: Find wallet", func(t *testing.T) {
		test := "f4c863ec-0300-495d-852d-c115e197390b"
		mockRow := sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow("f4c863ec-0300-495d-852d-c115e197390b", 1000, time.Now(), time.Now())

		db, mock, _ := sqlmock.New()
		defer db.Close()

		q := "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1"
		qExp := mock.ExpectQuery(q).WithArgs(test)
		qExp.WillReturnRows(mockRow)

		_, err := service.GetWalletsService(db, test)

		if err != nil {
			t.Errorf("TestGetWalletsService: got %v, want %v", err, nil)
		}
	})
}

func TestHandleOperationService(t *testing.T) {
	t.Run("Test 1: Deposit success", func(t *testing.T) {
		testWalletID := "f4c863ec-0300-495d-852d-c115e197390b"
		startBalance := 1000
		amount := 500

		db, mock, _ := sqlmock.New()
		defer db.Close()

		expectTxWithBalance(mock, testWalletID, startBalance, amount, nil)

		err := service.HandleOperationService(db, testWalletID, "DEPOSIT", amount)
		if err != nil {
			t.Errorf("HandleOperationService (DEPOSIT): got %v, want nil", err)
		}
	})

	t.Run("Test 2: Withdraw success", func(t *testing.T) {
		testWalletID := "f4c863ec-0300-495d-852d-c115e197390b"
		startBalance := 1500
		amount := 300

		db, mock, _ := sqlmock.New()
		defer db.Close()

		expectTxWithBalance(mock, testWalletID, startBalance, -amount, nil)

		err := service.HandleOperationService(db, testWalletID, "WITHDRAW", amount)
		if err != nil {
			t.Errorf("HandleOperationService (WITHDRAW): got %v, want nil", err)
		}
	})

	t.Run("Test 3: Withdraw causes negative balance", func(t *testing.T) {
		testWalletID := "f4c863ec-0300-495d-852d-c115e197390b"
		startBalance := 1000
		amount := 1500

		db, mock, _ := sqlmock.New()
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id, balance.*FOR UPDATE").
			WithArgs(testWalletID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
				AddRow(testWalletID, startBalance, time.Now(), time.Now()))
		mock.ExpectRollback()

		err := service.HandleOperationService(db, testWalletID, "WITHDRAW", amount)
		if !errors.Is(err, utils.ErrNegativeBalance) && !errors.Is(err, utils.ErrInvalidAmount) {
			t.Errorf("HandleOperationService: got %v, want negative balance error", err)
		}
	})

	t.Run("Test 4: Repository exec error", func(t *testing.T) {
		testWalletID := "f4c863ec-0300-495d-852d-c115e197390b"
		startBalance := 1000
		amount := 200

		db, mock, _ := sqlmock.New()
		defer db.Close()

		expectTxWithBalance(mock, testWalletID, startBalance, amount, sql.ErrConnDone)

		err := service.HandleOperationService(db, testWalletID, "DEPOSIT", amount)
		if err == nil {
			t.Error("HandleOperationService: expected error, got nil")
		}
	})
}
