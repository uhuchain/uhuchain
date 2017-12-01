//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package handler

import (
	"fmt"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/core/models"
	"github.com/uhuchain/uhuchain-api/interfaces/rest/operations/status"
)

// TODO Move GetStatusFailed to code generator / api definition

// HandleStatus implements the handling of the status endpoint
func (handler *RequestHandler) HandleStatus(params status.GetStatusParams) middleware.Responder {
	payload := &models.APIResponse{}
	ledgerStatus, err := handler.uhuClient.GetBlockchainInfo()
	if err != nil {
		log.Fatalf("Failed to get blockchain info. %s", err)
		resError := status.NewGetStatusInternalServerError()
		payload = models.NewErrorResponse(1200, fmt.Sprintf("Unable to get the uhuchain ledger status. %s", err))
		return resError.WithPayload(payload)
	}
	res := status.NewGetStatusOK()
	payload = models.NewMessageResponse(1000, fmt.Sprintf("Uhuchain car ledger API is alive. Ledger status: %s", ledgerStatus))
	res.WithPayload(payload)
	return res
}
