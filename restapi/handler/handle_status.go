package handler

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/ledger"
	"github.com/uhuchain/uhuchain-api/models"
	"github.com/uhuchain/uhuchain-api/restapi/operations/status"
)

// HandleStatus implements the handling of the status endpoint
func HandleStatus(params status.GetStatusParams) middleware.Responder {
	c := ledger.Client{}
	c.Init()
	currentBlock := c.GetCurrentBlock()
	res := status.NewGetStatusOK()
	payload := models.APIResponse{
		Code:    1000,
		Message: fmt.Sprintf("Uhuchain car ledger API is alive. Current block %s", currentBlock),
		Type:    "message",
	}
	res.WithPayload(&payload)
	return res
}
