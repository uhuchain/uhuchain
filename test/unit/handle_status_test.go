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
	"github.com/uhuchain/uhuchain-api/restapi/operations/status"
)

func TestHandleStatus(t *testing.T) {
	requestHandler := handler.NewRequestHandler(&ClientMock{})

	type args struct {
		params status.GetStatusParams
	}
	tests := []struct {
		name string
		args args
		want middleware.Responder
	}{
		{
			name: "Get status",
			args: args{
				status.GetStatusParams{},
			},
			want: status.NewGetStatusOK().WithPayload(models.NewMessageResponse(1000,
				"Uhuchain car ledger API is alive. Ledger status: Current block height:9 currentBlockHash:\"&\\021\\005_\\021\\325\\343T\\337\\213S\\206b?h\\2252\\010Y\\t\\331\\236\\213\\372\\244:V\\377\\324\\343e\\332\" previousBlockHash:\"\\343\\202O\\261\\367t\\310\\267\\242\\223\\177$[\\260\\022C={\\247\\266\\341\\207\\016\\r\\020\\337<\\211o++a\" ")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestHandler.HandleStatus(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
