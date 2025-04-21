package repositories

import (
	"JavaCode/internal/models"
	"JavaCode/utils"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

// GetWalletByUUID retrieves a wallet by UUID.
//
// Parameters:
//   - db: DB connection or transaction
//   - walletUUID: wallet identifier
//
// Returns:
//   - the wallet if found
//   - utils.ErrWalletNotFound if not found
//   - any other error on failure
func GetWalletByUUID(db Querier, walletUUID string) (*models.Wallet, error) {
	var wallet models.Wallet
	const query = "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = $1"
	err := db.QueryRow(query, walletUUID).Scan(&wallet.Id, &wallet.Balance, &wallet.CreatedTime, &wallet.UpdatedTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrWalletNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

// GetWalletForUpdate retrieves and locks a wallet by UUID.
//
// Parameters:
//   - db: transactional context (e.g., *sql.Tx)
//   - walletUUID: wallet identifier
//
// Returns:
//   - the wallet if found
//   - utils.ErrWalletNotFound if not found
//   - any other error on failure
func GetWalletForUpdate(db Querier, walletUUID string) (*models.Wallet, error) {
	var wallet models.Wallet
	query := "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = $1 FOR UPDATE"
	err := db.QueryRow(query, walletUUID).
		Scan(&wallet.Id, &wallet.Balance, &wallet.CreatedTime, &wallet.UpdatedTime)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrWalletNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

// ChainBalance updates the wallet's balance by delta.
//
// Parameters:
//   - db: DB connection or transaction
//   - walletUUID: wallet identifier
//   - delta: amount to add/subtract
//
// Returns:
//   - nil if successful
//   - utils.ErrNegativeBalance if balance goes below zero
//   - utils.ErrWalletNotFound if wallet doesn't exist
//   - any other error on failure
func ChainBalance(db Querier, walletUUID string, delta int) error {
	const query = "UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2"
	result, err := db.Exec(query, delta, walletUUID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Constraint == "wallets_balance_check" {
			return utils.ErrNegativeBalance
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return utils.ErrWalletNotFound
	}

	return nil
}
