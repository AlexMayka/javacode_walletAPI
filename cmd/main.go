// Package main is the entry point of the Wallet API service.
//
// This service provides a simple REST API for managing wallet operations
// such as balance inquiry, deposit, and withdrawal.
//
// Features:
//   - PostgreSQL database connection
//   - Config loading from environment
//   - REST API with Gin framework
//   - Middleware-based structured logging
//   - Swagger documentation support
//
// Endpoints:
//   - GET    /api/v1/wallets/{wallet_uuid} — get wallet balance
//   - POST   /api/v1/wallet                — perform deposit or withdrawal

// @title Wallet API
// @version 1.0
// @description API for wallet operation
// @BasePath /api/v1
package main

import (
	"JavaCode/config"
	"JavaCode/internal/routes"
	"JavaCode/pkg/db"
	"JavaCode/utils"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	utils.InitLogger()
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.Db,
	)

	utils.Logger.Infof("Connect to DataBase: %v", dsn)
	dbConn, err := db.InitDB(dsn, cfg.Db.Driver)
	if err != nil {
		utils.Logger.Fatalf("Failed to init DB: %v", err)
	}
	utils.Logger.Infof("Luck connect to DataBase: %v", dsn)
	defer dbConn.Close()

	dbConn.SetMaxOpenConns(100)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxLifetime(time.Hour)

	router := routes.SetupRouter(dbConn)

	addr := cfg.Host.ServerHost + ":" + cfg.Host.ServerPort

	utils.Logger.Infof("Start listing server: %v", addr)
	router.Run(addr)

	utils.Logger.Infof("Finish listing server: %v", addr)
}
