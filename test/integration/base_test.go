/*
Copyright Uhuchain. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"testing"

	"github.com/uhuchain/uhuchain-api/ledger"
)

var setup ledger.BaseSetupImpl

var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}

func init() {
	log.Println("Initialize uhuchain ledger client ")
	setup = ledger.BaseSetupImpl{
		ConfigFile:      "../client-config/config.yaml",
		ChannelID:       "car-ledger",
		OrgID:           "InsuranceA",
		OrdererOrgID:    "UhuchainOrderer",
		ChannelConfig:   "../uhuchain-network-dev/channel-artifact/channel.tx",
		ConnectEventHub: true,
		UserID:          "User1",
		AdminUserID:     "Admin",
	}
}

func TestClient(t *testing.T) {
	err := setup.Initialize()
	if err != nil {
		t.Fatalf("Failed to init client setup: %s", err)
	}
}
