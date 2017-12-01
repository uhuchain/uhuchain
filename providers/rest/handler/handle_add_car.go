//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package handler

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/core/models"
	"github.com/uhuchain/uhuchain-api/providers/rest/operations/car"
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

	err = handler.uhuClient.WriteToLedger(handler.carChainCode, strconv.FormatInt(newCar.ID, 10), carValue)
	if err != nil {
		res := car.NewAddCarInternalServerError()
		message := models.NewErrorResponse(500, err.Error())
		return res.WithPayload(message)
	}
	return car.NewAddCarCreated()
}
