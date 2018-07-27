// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int64  `json:"status"`
}

func handleErrorResponse(response *http.Response) error {
	var errorResponse ErrorResponse
	if err := unmarshalResponseBody(response, &errorResponse); err != nil {
		return err
	}
	return errors.New(errorResponse.Message)
}

func handleUnknownResponse(response *http.Response) error {
	var errStr = fmt.Sprintf("Unknown http status %d\n", response.StatusCode)
	return errors.New(errStr)
}
