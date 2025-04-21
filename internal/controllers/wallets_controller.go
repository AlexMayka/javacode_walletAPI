package controllers

import (
	"JavaCode/internal/models"
	"JavaCode/internal/service"
	"JavaCode/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// GetBalanceHandler godoc
// @Summary  Get Balance
// @Description  Return balance by UUID
// @Tags     wallet
// @Param    WALLET_UUID path string true "UUID wallet"
// @Success  200 {object} models.BalanceResponse
// @Failure  400 {object} utils.ErrorResponse
// @Failure  404 {object} utils.ErrorResponse
// @Router   /wallets/{WALLET_UUID} [get]
func (controller *Controller) GetBalanceHandler(c *gin.Context) {
	walletUUID := c.Param("WALLET_UUID")

	if _, err := uuid.Parse(walletUUID); err != nil {
		utils.Logger.WithError(err).Warn("Invalid uuid")
		utils.HandleError(c, utils.ErrInvalidRequest)
		return
	}

	wallet, err := service.GetWalletsService(controller.DB, walletUUID)
	if err != nil {
		utils.Logger.WithError(err).Warn("service GetWalletService failed")
		utils.HandleError(c, err)
		return
	}

	result := models.BalanceResponse{Uuid: wallet.Id, Balance: wallet.Balance}
	c.JSON(http.StatusOK, result)
}

// WalletOperationHandler godoc
// @Summary      Perform a wallet operation
// @Description  Deposit funds to, or withdraw funds from, a wallet.
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        request  body      models.WalletOperationRequest  true  "Operation parameters"
// @Success      200      {object}  map[string]string              "Operation successful"
// @Failure      400      {object}  utils.ErrorResponse            "Invalid request / negative amount"
// @Failure      404      {object}  utils.ErrorResponse            "Wallet not found"
// @Failure      500      {object}  utils.ErrorResponse            "Internal server error"
// @Router       /wallet [post]
func (controller *Controller) WalletOperationHandler(c *gin.Context) {
	var request models.WalletOperationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleError(c, utils.ErrInvalidRequest)
		utils.Logger.WithError(err).Warn("bad JSON body")
		return
	}

	if request.Amount <= 0 {
		utils.Logger.Warn("amount must be greater than zero")
		utils.HandleError(c, utils.ErrNegativeBalance)
		return
	}

	if err := ValidateUUID(request.WalletID); err != nil {
		utils.Logger.WithError(err).Warn("invalid UUID")
		utils.HandleError(c, err)
		return
	}

	if err := ValidateOperationType(request.OperationType); err != nil {
		utils.Logger.WithError(err).Warn("operation failed")
		utils.HandleError(c, err)
		return
	}

	err := service.HandleOperationService(controller.DB, request.WalletID, request.OperationType, request.Amount)
	if err != nil {
		utils.Logger.WithError(err).Warn("service Handle Operation failed")
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operation successful"})
}

// ValidateUUID checks if the given string is a valid UUID format.
//
// Returns:
//   - nil if the UUID is valid;
//   - utils.ErrInvalidRequest if the format is invalid.
func ValidateUUID(walletUUID string) error {
	if _, err := uuid.Parse(walletUUID); err != nil {
		return utils.ErrInvalidRequest
	}
	return nil
}

// ValidateOperationType checks if the operation type is valid.
//
// Allowed values are:
//   - service.DEPOSIT
//   - service.WITHDRAW
//
// Returns:
//   - nil if the operation type is valid;
//   - utils.ErrInvalidRequest if it's not one of the allowed values.
func ValidateOperationType(operationType string) error {
	if operationType != service.DEPOSIT && operationType != service.WITHDRAW {
		return utils.ErrInvalidRequest
	}
	return nil
}
