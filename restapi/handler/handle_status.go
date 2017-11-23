package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/ledger"
	"github.com/uhuchain/uhuchain-api/models"
	"github.com/uhuchain/uhuchain-api/restapi/operations/status"
)

// TODO Move GetStatusFailed to code generator / api definition

// GetStatusFailed message
type GetStatusFailed struct {
	Payload *models.APIResponse `json:"body,omitempty"`
}

// WriteResponse to the client
func (m *GetStatusFailed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(500)
	if m.Payload != nil {
		payload := m.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// HandleStatus implements the handling of the status endpoint
func HandleStatus(params status.GetStatusParams) middleware.Responder {
	payload := &models.APIResponse{}
	ledgerStatus, err := ledger.UhuClient.Setup.Channel.QueryInfo()
	if err != nil {
		log.Fatalf("Failed to get current hashblock. %s", err)
		*payload = models.APIResponse{
			Code:    1200,
			Message: fmt.Sprintf("Uhuchain car ledger API is alive. Current block %s", ledgerStatus.String()),
			Type:    "error",
		}
		return &GetStatusFailed{Payload: payload}
	}
	res := status.NewGetStatusOK()
	*payload = models.APIResponse{
		Code:    1000,
		Message: fmt.Sprintf("Uhuchain car ledger API is alive. Current block %s", ledgerStatus.String()),
		Type:    "message",
	}
	res.WithPayload(payload)
	return res
}
