/*
Copyright Uhuchain. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/uhuchain/uhuchain-api/ledger"
)

var setup ledger.BaseSetupImpl

var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}

func init() {
	log.SetOutput(os.Stdout)
	log.Println("Initialize uhuchain ledger client ")
	setup = ledger.BaseSetupImpl{
		ConfigFile:      os.Getenv("UHU_CONFIG"),
		ChannelID:       os.Getenv("UHU_CHANNELNAME"),
		OrgID:           os.Getenv("UHU_ORG"),
		OrdererOrgID:    os.Getenv("UHU_ORDERER"),
		ChannelConfig:   os.Getenv("UHU_CHANNELCONFIG"),
		ConnectEventHub: true,
		UserID:          os.Getenv("UHU_USERID"),
		AdminUserID:     os.Getenv("UHU_ADMINID"),
	}
}

func TestClient(t *testing.T) {
	err := setup.Initialize()
	if err != nil {
		t.Fatalf("Failed to init client setup: %s", err)
	}

	blockchaininfo, err := setup.Channel.QueryInfo()
	if err != nil {
		log.Fatalf("Failed to get current hashblock. %s", err)
	}
	block, err := setup.Channel.QueryBlockByHash(blockchaininfo.CurrentBlockHash)
	if err != nil {
		log.Fatalf("QueryBlockByHash return error: %v", err)
	}

	t.Logf("Current block: %s", block.String())
	t.Logf("Total number of blocks: %d", blockchaininfo.GetHeight())
	t.Logf("Peers %s", setup.Channel.Peers())
	t.Logf("Blockchain info %s", blockchaininfo.String())

	t.Logf("Primary peers %s", setup.Channel.PrimaryPeer().URL())

	testQuery("automotive", setup.ChannelClient, t)

}

func testQuery(ccID string, chClient apitxn.ChannelClient, t *testing.T) {
	t.Log("Query ledger")
	var queryArgs = [][]byte{[]byte("b")}
	result, err := chClient.Query(apitxn.QueryRequest{ChaincodeID: ccID, Fcn: "query", Args: queryArgs})
	if err != nil {
		t.Fatalf("Failed to invoke example cc: %s", err)
	}
	t.Logf("Query result: %s", result)
}
