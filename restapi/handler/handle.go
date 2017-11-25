package handler

import (
	"github.com/uhuchain/uhuchain-api/ledger"
	"github.com/uhuchain/uhuchain-api/models"
)

var uhuClient ledger.Client

var carChainCode = "cc_car"

// SetLedgerClient sets the ledger client for the handler package
func SetLedgerClient(uhuchainClient ledger.Client) {
	uhuClient = uhuchainClient
}

func createMessage(code int32, message string, messageType string) *models.APIResponse {
	return &models.APIResponse{
		Code:    code,
		Message: message,
		Type:    messageType,
	}
}
