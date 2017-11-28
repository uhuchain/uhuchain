//
// Copyright Uhuchain All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package models

func newApiResponse(code int32, message string, messageType string) *APIResponse {
	return &APIResponse{
		Code:    code,
		Message: message,
		Type:    messageType,
	}
}

func NewMessageResponse(code int32, message string) *APIResponse {
	return newApiResponse(code, message, "message")
}

func NewErrorResponse(code int32, message string) *APIResponse {
	return newApiResponse(code, message, "error")
}
