package handler

import (
	"github.com/uhuchain/uhuchain-api/ledger"
)

var uhuClient ledger.Client

var carChainCode = "automotive"

// SetLedgerClient sets the ledger client for the handler package
func SetLedgerClient(uhuchainClient ledger.Client) {
	uhuClient = uhuchainClient
}
