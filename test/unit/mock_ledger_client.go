package unit

import (
	"errors"
	"fmt"

	"github.com/uhuchain/uhuchain-core/models"
)

// ClientMock mocks the hlf client
type ClientMock struct {
	QueryResponse string
}

//GetBlockchainInfo mock
func (mock *ClientMock) GetBlockchainInfo() (string, error) {
	return "Current block height:9 currentBlockHash:\"&\\021\\005_\\021\\325\\343T\\337\\213S\\206b?h\\2252\\010Y\\t\\331\\236\\213\\372\\244:V\\377\\324\\343e\\332\" previousBlockHash:\"\\343\\202O\\261\\367t\\310\\267\\242\\223\\177$[\\260\\022C={\\247\\266\\341\\207\\016\\r\\020\\337<\\211o++a\" ", nil
}

//QueryLedger mock. Use id 12345 to get a car object and everything else to get an error.
func (mock *ClientMock) QueryLedger(ccID string, fnc string, id string) ([]byte, error) {
	if id == "12345" {
		return []byte(mock.QueryResponse), nil
	}
	errorMessage := fmt.Sprintf("SendTransactionProposal failed: Transaction processor (peer0.insurancea.uhuchain.com:7051) returned error 'rpc error: code = Unknown desc = chaincode error (status: 500, message: {\"Error\":\"Nil amount for %s\"})' for txID '17d11ca51f7512a5dc4a7db84a7f0ef75d057e9acdad0ff5bebcf68d4dd9f63b'", id)
	return []byte(""), errors.New(errorMessage)
}

//Invoke mocks the functions on the chaincode
func (mock *ClientMock) Invoke(ccID string, fnc string, args [][]byte) error {
	if fnc == "saveCar" {
		car := models.Car{}
		err := car.UnmarshalBinary(args[0])
		if err != nil {
			return err
		}
		if car.ID == 12345 {
			return nil
		}
		return errors.New("unable to write onto the ledger")
	}

	return errors.New("not implemented yet")
}

//Init mock
func (mock *ClientMock) Init() {
}
