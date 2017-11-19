/*
Copyright Uhuchain. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/uhuchain/uhu-hlf-client/config"
	"github.com/uhuchain/uhu-hlf-client/log"
)

var configFile = "../fixtures/config/config.yaml"
var orgID = "Org1"
var ordererOrg = "ordererorg"
var orgUser = "Admin"
var ordererUser = "Admin"

var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}

func init() {
	log.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func TestChannelOperations(t *testing.T) {
	client := &config.ClientConfig{
		ConfigFile: configFile,
	}

	err := client.Initialize(orgID, orgUser)
	if err != nil {
		t.Errorf("Failed to init blockchain client. Message: %s", err)
	}

	channel, err := client.CreateChannel("mychannel", ordererOrg, ordererUser, []string{orgID},
		"../fixtures/channel/mychannel.tx")

	if err != nil {
		t.Errorf("Failed to create channel. Message: %s", err)
	}

	eventHub, err := client.GetEventHubForOrg(orgID)
	if err != nil {
		t.Errorf("Failed to get eventhub for org %s. Message: %s", orgID, err)
	}

	// Always connect to the eventhub prior to using it, otherwise functions using
	// the eventhub may fail
	err = eventHub.Connect()
	if err != nil {
		t.Errorf("Failed to connect to eventhub")
	}

	err = client.SetupCC(channel, "testcode", GetDeployPath(), "github.com/example_cc",
		"v0", nil, initArgs, eventHub)
	if err != nil {
		t.Errorf("Failed to install chaincode. Message: %s", err)
	}
	t.Logf("Query chaincode on channel %s ...", channel.Name())
	var queryArgs = [][]byte{[]byte("query"), []byte("b")}
	value, err := client.Query("mychannel", "User1", "testcode", queryArgs)
	//value, err := client.QueryChannel(channel, "testcode", queryArgs, eventHub)
	if err != nil {
		t.Errorf("Failed to query! Message: %s", err)
	}
	t.Logf("Value return from query is %s", value)
	if string(value) != "200" {
		t.Errorf("Got result %s from query. Should be 200", string(value))
	}
}

// GetDeployPath ..
func GetDeployPath() string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, "../fixtures/testdata")
}
