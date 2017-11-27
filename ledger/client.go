/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"log"
	"os"
	"time"

	"github.com/uhuchain/uhuchain-api/models"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// Client defines the interface for a client that interacts with the ledger
type Client interface {
	GetBlockchainInfo() (string, error)
	QueryLedger(string, string) ([]byte, error)
	WriteToLedger(string, string, []byte) error
	Init()
}

// FabricClient wraps the hlf fabric sdk client
type FabricClient struct {
	//Setup client
	setup BaseSetupImpl
}

type TestTxFilter struct {
	err          error
	errResponses error
}

func (tf *TestTxFilter) ProcessTxProposalResponse(txProposalResponse []*apitxn.TransactionProposalResponse) ([]*apitxn.TransactionProposalResponse, error) {
	if tf.err != nil {
		return nil, tf.err
	}

	var newResponses []*apitxn.TransactionProposalResponse

	if tf.errResponses != nil {
		// 404 will cause transaction commit error
		txProposalResponse[0].ProposalResponse.Response.Status = 404
	}

	newResponses = append(newResponses, txProposalResponse[0])
	return newResponses, nil
}

// GetBlockchainInfo returns the current block of the channel
func (client *FabricClient) GetBlockchainInfo() (string, error) {
	log.Println("Getting current block")
	blockchaininfo, err := client.setup.Channel.QueryInfo()

	if err != nil {
		log.Fatalf("Failed to get current blcokchaininfo. %s", err)
		return "", err
	}
	return blockchaininfo.String(), nil
}

// QueryLedger performs a basic query operation on the ledger
func (client *FabricClient) QueryLedger(ccID string, id string) ([]byte, error) {
	log.Print("Query ledger")
	var queryArgs = [][]byte{[]byte(id)}
	result, err := client.setup.ChannelClient.Query(apitxn.QueryRequest{ChaincodeID: ccID, Fcn: "query", Args: queryArgs})
	if err != nil {
		log.Fatalf("Failed to invoke example cc: %s", err)
		return nil, err
	}
	log.Printf("Query result: %s", result)
	return result, nil
}

// WriteToLedger writes an object onto the ledger
func (client *FabricClient) WriteToLedger(ccID string, carID string, value []byte) error {
	txNotifier := make(chan apitxn.ExecuteTxResponse)
	txFilter := &TestTxFilter{}
	txOpts := apitxn.ExecuteTxOpts{Notifier: txNotifier, TxFilter: txFilter}

	var newCarArg = [][]byte{[]byte(carID), value}

	_, err := client.setup.ChannelClient.ExecuteTxWithOpts(apitxn.ExecuteTxRequest{ChaincodeID: ccID, Fcn: "write", Args: newCarArg}, txOpts)
	if err != nil {
		log.Fatalf("Failed to move funds: %s", err)
		return err
	}

	select {
	case response := <-txNotifier:
		if response.Error != nil {
			log.Fatalf("ExecuteTx returned error: %s", response.Error)
			return err
		}
		if response.TxValidationCode != pb.TxValidationCode_VALID {
			log.Fatalf("Expecting TxValidationCode to be TxValidationCode_VALID but received: %s", response.TxValidationCode)
			return err
		}
	case <-time.After(time.Second * 20):
		log.Fatalf("ExecuteTx timed out")
		return err
	}
	return nil
}

// Init set the client
func (client *FabricClient) Init() {
	client.setup = BaseSetupImpl{
		ConfigFile:      os.Getenv("UHU_CONFIG"),
		ChannelID:       os.Getenv("UHU_CHANNELNAME"),
		OrgID:           os.Getenv("UHU_ORG"),
		OrdererOrgID:    os.Getenv("UHU_ORDERER"),
		ChannelConfig:   os.Getenv("UHU_CHANNELCONFIG"),
		ConnectEventHub: true,
		UserID:          os.Getenv("UHU_USERID"),
		AdminUserID:     os.Getenv("UHU_ADMINID"),
	}
	err := client.setup.Initialize()
	if err != nil {
		log.Fatalf("Failed to init client setup: %s", err)
		os.Exit(1)
	}
	log.Println("Uhuchain client initilized successfully.")
}

// ClientMock mocks the hlf client
type ClientMock struct {
}

func (mock *ClientMock) GetBlockchainInfo() (string, error) {
	res := models.APIResponse{
		Code:    1000,
		Message: "Uhuchain car ledger API is alive. Current block height:9 currentBlockHash:\"&\\021\\005_\\021\\325\\343T\\337\\213S\\206b?h\\2252\\010Y\\t\\331\\236\\213\\372\\244:V\\377\\324\\343e\\332\" previousBlockHash:\"\\343\\202O\\261\\367t\\310\\267\\242\\223\\177$[\\260\\022C={\\247\\266\\341\\207\\016\\r\\020\\337<\\211o++a\" ",
		Type:    "message",
	}
	status, err := res.MarshalBinary()
	return string(status), err
}

func (mock *ClientMock) QueryLedger(string, string) ([]byte, error) {
	panic("not implemented")
}

func (mock *ClientMock) WriteToLedger(string, string, []byte) error {
	panic("not implemented")
}

func (mock *ClientMock) Init() {
	panic("not implemented")
}
