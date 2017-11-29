//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package unit

import (
	"reflect"
	"testing"

	"github.com/uhuchain/uhuchain-api/models"

	"github.com/go-openapi/runtime/middleware"
	"github.com/uhuchain/uhuchain-api/restapi/handler"
	"github.com/uhuchain/uhuchain-api/restapi/operations/car"
)

func TestHandleGetCar(t *testing.T) {
	payload := `{
		"brand": "Volkswagen",
		"id": 12345,
		"model": "Sharan GTI",
		"policies": [
			{
				"claims": [
					{
						"date": "2016-11-01",
						"description": "Something bad happend",
						"id": 12345
					}
				],
				"endDate": "2017-09-01",
				"id": 12345,
				"insuranceId": 12345,
				"insuranceName": "Zurich Insurance Group",
				"startDate": "2016-09-01"
			}
		],
		"vehicleId": "THK34SDM6A2D34"
	}`
	mockClient := &ClientMock{
		QueryResponse: payload,
	}
	requestHandler := handler.NewRequestHandler(mockClient)

	carResult := models.Car{}
	carResult.UnmarshalBinary([]byte(payload))

	getCarParamsOk := car.NewGetCarParams()
	getCarParamsOk.ID = 12345

	getCarParamsFail := car.NewGetCarParams()
	getCarParamsFail.ID = 22222

	type args struct {
		params car.GetCarParams
	}
	tests := []struct {
		name string
		args args
		want middleware.Responder
	}{
		{
			name: "Get car successfully",
			args: args{
				getCarParamsOk,
			},
			want: car.NewGetCarOK().WithPayload(&carResult),
		},
		{
			name: "Car not found",
			args: args{
				getCarParamsFail,
			},
			want: car.NewGetCarNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestHandler.HandleGetCar(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleGetCar() = %v, want %v", got, tt.want)
			}
		})
	}
}
