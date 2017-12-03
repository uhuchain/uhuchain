// +build !test

/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package hlf

import (
	"log"
	"os"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

func init() {
	envVariables := [7]string{"UHU_CONFIG",
		"UHU_CHANNELNAME",
		"UHU_ORG", "UHU_ORDERER",
		"UHU_CHANNELCONFIG",
		"UHU_USERID",
		"UHU_ADMINID",
	}
	for i := range envVariables {
		v := envVariables[i]
		if _, exists := os.LookupEnv(v); exists == false {
			log.Fatalf("Failed to init uhuchain server. $%s not found.", v)
			os.Exit(1)
		}
	}
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
func (client *FabricClient) QueryLedger(ccID string, fcn string, id string) ([]byte, error) {
	log.Print("Query ledger")
	var queryArgs = [][]byte{[]byte(id)}
	result, err := client.setup.ChannelClient.Query(apitxn.QueryRequest{ChaincodeID: ccID, Fcn: fcn, Args: queryArgs})
	if err != nil {
		log.Fatalf("Failed to invoke example cc: %s", err)
		return nil, err
	}
	log.Printf("Query result: %s", result)
	return result, nil
}

// Invoke a chaincode function on the ledger
func (client *FabricClient) Invoke(ccID string, fcn string, args [][]byte) error {
	log.Printf("Invoking function %s", fcn)
	txNotifier := make(chan apitxn.ExecuteTxResponse)
	txFilter := &TestTxFilter{}
	txOpts := apitxn.ExecuteTxOpts{Notifier: txNotifier, TxFilter: txFilter}

	_, err := client.setup.ChannelClient.ExecuteTxWithOpts(apitxn.ExecuteTxRequest{ChaincodeID: ccID, Fcn: fcn, Args: args}, txOpts)
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
