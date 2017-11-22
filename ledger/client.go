/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledger

import (
	"log"
	"os"
)

// Client wraps the hlf fabric sdk client
type Client struct {
	//Setup client
	Setup BaseSetupImpl
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
