package handler

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/models"
	"github.com/uhuchain/uhuchain-api/restapi/operations/car"
)

// HandleGetCar gets the car from the ledger
func HandleGetCar(params car.GetCarParams) middleware.Responder {
	carResult := &models.Car{}
	ccID := "cc_car"
	res := car.NewGetCarOK()
	result, err := uhuClient.QueryLedger(ccID, strconv.FormatInt(params.ID, 10))
	if err != nil {
		res := car.NewGetCarInternalServerError()
		message := createMessage(500, err.Error(), "error")
		return res.WithPayload(message)
	}

	carResult.UnmarshalBinary(result)
	res.WithPayload(carResult)
	return res
}
