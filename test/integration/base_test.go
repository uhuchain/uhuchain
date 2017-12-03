/*
Copyright Uhuchain. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"os"
	"testing"

	"github.com/uhuchain/uhuchain-api/providers/hlf"
)

var setup hlf.BaseSetupImpl

func init() {
	log.SetOutput(os.Stdout)
	log.Println("Initialize uhuchain ledger client ")
	setup = hlf.BaseSetupImpl{
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

// Test code that actually interacts with a running fabric framework
// TODO Write full integration test

func TestClient(t *testing.T) {

	client := hlf.FabricClient{}
	client.Init()
	result, err := client.GetBlockchainInfo()
	if err != nil {
		t.Fatalf("Failed to call ledger - %s", err)
	}
	t.Logf("Got %s", result)
}
