package handler

import (
	"fmt"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/models"
	"github.com/uhuchain/uhuchain-api/restapi/operations/status"
)

// TODO Move GetStatusFailed to code generator / api definition

// HandleStatus implements the handling of the status endpoint
func HandleStatus(params status.GetStatusParams) middleware.Responder {
	payload := &models.APIResponse{}
	ledgerStatus, err := uhuClient.GetBlockchainInfo()
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
