package service

import (
	"JavaCode/internal/models"
	"JavaCode/internal/repositories"
	"JavaCode/utils"
	"database/sql"
	"errors"
	"fmt"
)

const (
	DEPOSIT  = "DEPOSIT"
	WITHDRAW = "WITHDRAW"
)

// GetWalletsService retrieves a wallet by UUID.
//
// It returns:
//   - the wallet if found;
//   - utils.ErrWalletNotFound if the wallet does not exist;
//   - any other error from the repository layer.
func GetWalletsService(db *sql.DB, walletUUID string) (*models.Wallet, error) {
	wallet, err := repositories.GetWalletByUUID(db, walletUUID)
	if err != nil {
		if errors.Is(err, utils.ErrWalletNotFound) {
			return nil, utils.ErrWalletNotFound
		}
		return nil, utils.ErrDatabase
	}
	return wallet, nil
}

// HandleOperationService processes a deposit or withdrawal operation on a wallet.
//
// It calculates the delta (positive or negative) based on the operation type,
// and applies the change via the repository layer.
//
// Returns:
//   - nil on success;
//   - an error if the balance update fails.
func HandleOperationService(db *sql.DB, walletID, operationType string, amount int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx error: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	wallet, err := repositories.GetWalletForUpdate(tx, walletID)
	if err != nil {
		return err
	}

	var delta int
	if operationType == DEPOSIT {
		delta = amount
	} else if operationType == WITHDRAW {
		delta = -amount
	}

	newBalance := int(wallet.Balance) + delta
	if newBalance < 0 {
		return utils.ErrNegativeBalance
	}

	if err := repositories.ChainBalance(tx, walletID, delta); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return err
}
