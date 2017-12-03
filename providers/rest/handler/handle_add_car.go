//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/providers/rest/operations/car"
	"github.com/uhuchain/uhuchain-core/models"
)

// HandleAddCar adds a car from to ledger
func (handler *RequestHandler) HandleAddCar(params car.AddCarParams) middleware.Responder {
	newCar := params.Body

	carValue, err := newCar.MarshalBinary()
	if err != nil {
		res := car.NewAddCarBadRequest()
		message := models.NewErrorResponse(400, err.Error())
		return res.WithPayload(message)
	}

	var args = [][]byte{[]byte(carValue)}

	err = handler.uhuClient.Invoke(handler.carChainCode, "saveCar", args)
	if err != nil {
		res := car.NewAddCarInternalServerError()
		message := models.NewErrorResponse(500, err.Error())
		return res.WithPayload(message)
	}
	return car.NewAddCarCreated()
}
