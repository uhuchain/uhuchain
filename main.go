/*
Copyright Uhuchain. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"io/ioutil"
	"os"

	"github.com/uhuchain/uhu-hlf-client/config"
	"github.com/uhuchain/uhu-hlf-client/log"
)

func init() {
	log.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	log.Error.Print("sdsdf")
	log.Info.Println("### Starting hyperledger-client ###")
	setup := &config.ClientConfig{
		ConfigFile: "test/fixtures/config/config.yaml",
	}
	err := setup.Initialize("Org1", "Admin")
	if err != nil {
		log.Error.Printf("Failed to init blockchain client. Message: %s", err)
	}
	log.Info.Println("### Shutdown hyperledger-client ###")
}
