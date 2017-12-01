//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package handler

import "github.com/uhuchain/uhuchain-api/interfaces/hlf"

// RequestHandler to handle the uhuchain api requests
type RequestHandler struct {
	uhuClient    hlf.Client
	carChainCode string
}

// NewRequestHandler sets the ledger client for the handler package
func NewRequestHandler(uhuchainClient hlf.Client) RequestHandler {
	requestHandler := RequestHandler{}
	requestHandler.uhuClient = uhuchainClient
	requestHandler.uhuClient.Init()
	requestHandler.carChainCode = "automotive"
	return requestHandler
}
