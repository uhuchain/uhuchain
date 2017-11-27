package handler

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/restapi/operations/car"
)

// HandleAddCar adds a car from to ledger
func HandleAddCar(params car.AddCarParams) middleware.Responder {
	newCar := params.Body
	carValue, err := newCar.MarshalBinary()
	if err != nil {
		res := car.NewAddCarBadRequest()
		message := createMessage(400, err.Error(), "error")
		return res.WithPayload(message)
	}

	err = uhuClient.WriteToLedger(carChainCode, strconv.FormatInt(newCar.ID, 10), carValue)
	if err != nil {
		res := car.NewAddCarInternalServerError()
		message := createMessage(500, err.Error(), "error")
		return res.WithPayload(message)
	}
	return car.NewAddCarCreated()
}
