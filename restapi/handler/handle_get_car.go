//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package handler

import (
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/models"
	"github.com/uhuchain/uhuchain-api/restapi/operations/car"
)

// HandleGetCar gets the car from the ledger
func (handler *RequestHandler) HandleGetCar(params car.GetCarParams) middleware.Responder {
	carResult := &models.Car{}
	res := car.NewGetCarOK()
	result, err := handler.uhuClient.QueryLedger(handler.carChainCode, strconv.FormatInt(params.ID, 10))
	if err != nil {
		// TODO check if there is a better way to classify error responses
		if strings.Contains(err.Error(), "Nil amount for") {
			res := car.NewGetCarNotFound()
			return res
		}
		res := car.NewGetCarInternalServerError()
		message := models.NewErrorResponse(500, err.Error())
		return res.WithPayload(message)
	}

	carResult.UnmarshalBinary(result)
	res.WithPayload(carResult)
	return res
}
