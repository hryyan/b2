// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"net/http"
)

const AUTH_URL = "https://api.backblazeb2.com/b2api/v1/b2_authorize_account"

type AuthResponse struct {
	AccountId               string `json:"accountId"`
	AuthorizationToken      string `json:"authorizationToken"`
	ApiUrl                  string `json:"apiUrl"`
	DownloadUrl             string `json:"downloadUrl"`
	RecommendedPartSize     int64  `json:"recommendedPartSize"`
	AbsoluteMinimumPartSize int64  `json:"absoluteMinimumPartSize"`
}

// Auth your account
func (b *B2) Auth() error {
	request, err := http.NewRequest("GET", AUTH_URL, nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(b.AccountId, b.ApplicationKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		if err = unmarshalResponseBody(response, &b.auth); err != nil {
			return err
		}
		return nil
	case 401:
		return handleErrorResponse(response)
	default:
		return handleUnknownResponse(response)
	}
}
