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
	payload := models.APIResponse{}
	ledgerStatus, err := uhuClient.Setup.Channel.QueryInfo()
	if err != nil {
		log.Fatalf("Failed to get blockchain info. %s", err)
		resError := status.NewGetStatusInternalServerError()
		payload = models.APIResponse{
			Code:    1200,
			Message: fmt.Sprintf("Uhuchain car ledger API is alive. Current block %s", ledgerStatus.String()),
			Type:    "error",
		}
		return resError.WithPayload(&payload)
	}
	res := status.NewGetStatusOK()
	payload = models.APIResponse{
		Code:    1000,
		Message: fmt.Sprintf("Uhuchain car ledger API is alive. Current block %s", ledgerStatus.String()),
		Type:    "message",
	}
	res.WithPayload(&payload)
	return res
}
