package controllers_test

import (
	"JavaCode/internal/controllers"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func expectSuccessfulTx(mock sqlmock.Sqlmock, uuid string, delta int) {
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, balance.* FOR UPDATE").
		WithArgs(uuid).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
			AddRow(uuid, 1000, time.Now(), time.Now()))
	mock.ExpectExec("UPDATE wallets SET balance = balance.*").
		WithArgs(delta, uuid).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func TestValidateUUID(t *testing.T) {
	t.Run("Test 1: Test valid UUIDs", func(t *testing.T) {
		tests := []struct {
			name    string
			input   string
			wantErr error
		}{
			{"First", "f4c863ec-0300-495d-852d-c115e197390b", nil},
			{"Second", "2910473d-9d6a-45de-b87b-7a1b0d6540e6", nil},
			{"Third", "f7607b62-58a0-40f0-8fe7-bb54b01a516d", nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := controllers.ValidateUUID(tt.input)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				}
			})
		}
	})

	t.Run("Test 2: Test empty UUID", func(t *testing.T) {
		err := controllers.ValidateUUID("")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Test 3: Test invalid UUIDs", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"First", "f4c863ec-0300-495d-852d-97390b"},
			{"Second", "2910473d-9d6a--7a1b0d6540e6"},
			{"Third", "22222222323-58a0-40-8fe7-bb54b01a516d"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := controllers.ValidateUUID(tt.input)
				if err == nil {
					t.Errorf("expected error for input: %s", tt.input)
				}
			})
		}
	})
}

func TestValidateOperationType(t *testing.T) {
	t.Run("Test 1. Check DEPOSIT", func(t *testing.T) {
		test := "DEPOSIT"
		err := controllers.ValidateOperationType(test)
		if err != nil {
			t.Errorf("error = %v, got nil", err)
		}
	})

	t.Run("Test 2. Check WITHDRAW", func(t *testing.T) {
		test := "WITHDRAW"
		err := controllers.ValidateOperationType(test)
		if err != nil {
			t.Errorf("error = %v, got nil", err)
		}
	})

	t.Run("Test 3. Check another command", func(t *testing.T) {
		tests := []string{"ADD", "deposit", "withdraw", "DELETE", "SELECT"}
		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				err := controllers.ValidateOperationType(tt)
				if err == nil {
					t.Errorf("expected error, got error")
				}
			})
		}
	})
}

func TestController_GetBalanceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Test 1: Test Valid UUID", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			wantCode int
		}{
			{"Null UUID", "", http.StatusBadRequest},
			{"Invalid UUID", "1111-11111-11111-111111", http.StatusBadRequest},
			{"Valid UUID", "f4c863ec-0300-495d-852d-c115e197390b", http.StatusInternalServerError},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, _, _ := sqlmock.New()
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Params = gin.Params{{Key: "WALLET_UUID", Value: tt.input}}
				req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+tt.input, nil)
				c.Request = req

				ctrl := controllers.Controller{DB: db}
				ctrl.GetBalanceHandler(c)

				assert.Equal(t, tt.wantCode, w.Code)
			})
		}
	})

	t.Run("Test 2: Check Data", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			wantCode    int
			mockRows    *sqlmock.Rows
			mockErr     error
			expectQuery bool
		}{
			{
				name:        "Empty result",
				input:       "f4c863ec-0300-495d-852d-c115e197390b",
				wantCode:    http.StatusNotFound,
				mockErr:     sql.ErrNoRows,
				expectQuery: true,
			},
			{
				name:     "Ok",
				input:    "a1c122d7-fbc1-4ebb-bdd5-4ddb793c92bf",
				wantCode: http.StatusOK,
				mockRows: sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
					AddRow("a1c122d7-fbc1-4ebb-bdd5-4ddb793c92bf", 1000, time.Now(), time.Now()),
				expectQuery: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				defer db.Close()

				if tt.expectQuery {
					q := "SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1"
					qExp := mock.ExpectQuery(q).WithArgs(tt.input)

					if tt.mockErr != nil {
						qExp.WillReturnError(tt.mockErr)
					} else {
						qExp.WillReturnRows(tt.mockRows)
					}
				}

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				c.Params = gin.Params{{Key: "WALLET_UUID", Value: tt.input}}
				req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+tt.input, nil)
				c.Request = req

				ctrl := controllers.Controller{DB: db}
				ctrl.GetBalanceHandler(c)

				assert.Equal(t, tt.wantCode, w.Code)

				if tt.wantCode == http.StatusOK {
					assert.Contains(t, w.Body.String(), `"uuid"`)
					assert.Contains(t, w.Body.String(), tt.input)
					assert.Contains(t, w.Body.String(), `"balance"`)
				}
			})
		}
	})
}

func TestController_WalletOperationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Test 1: Test Body Request", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			wantCode int
		}{
			{
				name:     "Invalid JSON Body",
				input:    "{invalid-json",
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Missing Fields",
				input:    `{"walletId": "f4c863ec-0300-495d-852d-c115e197390b"}`,
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Negative Amount",
				input:    `{"walletId": "f4c863ec-0300-495d-852d-c115e197390b", "operationType": "DEPOSIT", "amount": -1000}`,
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Missing operationType",
				input:    `{"walletId": "f4c863ec-0300-495d-852d-c115e197390b", "amount": 1000}`,
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Valid Body",
				input:    `{"walletId": "f4c863ec-0300-495d-852d-c115e197390b", "operationType": "DEPOSIT", "amount": 1000}`,
				wantCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.wantCode == http.StatusOK {
					expectSuccessfulTx(mock, "f4c863ec-0300-495d-852d-c115e197390b", 1000)
				}

				body := strings.NewReader(tt.input)
				req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/", body)
				req.Header.Add("Content-Type", "application/json")

				c.Request = req

				ctrl := controllers.Controller{DB: db}
				ctrl.WalletOperationHandler(c)

				assert.Equal(t, tt.wantCode, w.Code)
			})
		}
	})

	t.Run("Test 2: Test valid UUIDs", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			wantCode int
		}{
			{
				name:     "Invalid UUID",
				input:    `{"walletId": "f4c8-030-495d-852d-c115e197390b", "operationType": "DEPOSIT", "amount": 1000}`,
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Null UUID",
				input:    `{"walletId": "", "operationType": "DEPOSIT", "amount": -1000}`,
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Valid UUID",
				input:    `{"walletId": "f4c863ec-0300-495d-852d-c115e197390b", "operationType": "DEPOSIT", "amount": 1000}`,
				wantCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				if tt.wantCode == http.StatusOK {
					expectSuccessfulTx(mock, "f4c863ec-0300-495d-852d-c115e197390b", 1000)
				}

				body := strings.NewReader(tt.input)
				req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet/", body)
				req.Header.Add("Content-Type", "application/json")

				c.Request = req

				ctrl := controllers.Controller{DB: db}
				ctrl.WalletOperationHandler(c)

				assert.Equal(t, tt.wantCode, w.Code)
			})
		}
	})

	t.Run("Test 3: Test operations", func(t *testing.T) {
		tests := []struct {
			name        string
			operation   string
			amount      int
			expectDelta int
			wantCode    int
		}{
			{"Deposit", "DEPOSIT", 1000, 1000, http.StatusOK},
			{"Withdraw", "WITHDRAW", 1000, -1000, http.StatusOK},
			{"Unknown", "TRANSFER", 1000, 0, http.StatusBadRequest},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				db, mock, _ := sqlmock.New()
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				// Ожидание транзакции только при валидной операции
				if tt.wantCode == http.StatusOK && (tt.operation == "DEPOSIT" || tt.operation == "WITHDRAW") {
					expectSuccessfulTx(mock, "f4c863ec-0300-495d-852d-c115e197390b", tt.expectDelta)
				}

				bodyStr := fmt.Sprintf(`{
				"walletId": "f4c863ec-0300-495d-852d-c115e197390b",
				"operationType": "%s",
				"amount": %d
			}`, tt.operation, tt.amount)

				req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(bodyStr))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req

				ctrl := controllers.Controller{DB: db}
				ctrl.WalletOperationHandler(c)

				assert.Equal(t, tt.wantCode, w.Code)
			})
		}
	})
}
