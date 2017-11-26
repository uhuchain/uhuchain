/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"log"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// Client wraps the hlf fabric sdk client
type Client struct {
	//Setup client
	Setup BaseSetupImpl
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

// GetCurrentBlock returns the current block of the channel
func (client *Client) GetCurrentBlock() string {
	log.Println("Getting current block")
	blockchaininfo, err := client.Setup.Channel.QueryInfo()

	if err != nil {
		log.Fatalf("Failed to get current hashblock. %s", err)
	}

	block, err := client.Setup.Channel.QueryBlockByHash(blockchaininfo.CurrentBlockHash)
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %v", err)
	}

	return string(block.String())
}

// QueryLedger performs a basic query operation on the ledger
func (client *Client) QueryLedger(ccID string, id string) ([]byte, error) {
	log.Print("Query ledger")
	var queryArgs = [][]byte{[]byte(id)}
	result, err := client.Setup.ChannelClient.Query(apitxn.QueryRequest{ChaincodeID: ccID, Fcn: "query", Args: queryArgs})
	if err != nil {
		log.Fatalf("Failed to invoke example cc: %s", err)
		return nil, err
	}
	log.Printf("Query result: %s", result)
	return result, nil
}

// WriteToLedger writes an object onto the ledger
func (client *Client) WriteToLedger(ccID string, carID string, value []byte) error {
	txNotifier := make(chan apitxn.ExecuteTxResponse)
	txFilter := &TestTxFilter{}
	txOpts := apitxn.ExecuteTxOpts{Notifier: txNotifier, TxFilter: txFilter}

	var newCarArg = [][]byte{[]byte(carID), value}

	_, err := client.Setup.ChannelClient.ExecuteTxWithOpts(apitxn.ExecuteTxRequest{ChaincodeID: ccID, Fcn: "write", Args: newCarArg}, txOpts)
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
func (client *Client) Init() {
	client.Setup = BaseSetupImpl{
		ConfigFile:      os.Getenv("UHU_CONFIG"),
		ChannelID:       os.Getenv("UHU_CHANNELNAME"),
		OrgID:           os.Getenv("UHU_ORG"),
		OrdererOrgID:    os.Getenv("UHU_ORDERER"),
		ChannelConfig:   os.Getenv("UHU_CHANNELCONFIG"),
		ConnectEventHub: true,
		UserID:          os.Getenv("UHU_USERID"),
		AdminUserID:     os.Getenv("UHU_ADMINID"),
	}
	err := client.Setup.Initialize()
	if err != nil {
		log.Fatalf("Failed to init client setup: %s", err)
		os.Exit(1)
	}
	log.Println("Uhuchain client initilized successfully.")
}
