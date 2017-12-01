package hlf

// Client defines the interface for a client that interacts with the ledger
type Client interface {
	GetBlockchainInfo() (string, error)
	QueryLedger(string, string) ([]byte, error)
	WriteToLedger(string, string, []byte) error
	Init()
}
