// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"fmt"
)

// CreateKey creates a new application key.
// See "b2_create_key" for an introduction:
// https://www.backblaze.com/b2/docs/b2_create_key.html
//
// Parameter capabilities, keyName are required.
// CreateKey return an ApplicationKey pointer and an error.
func (b *B2) CreateKey(capabilities []string, keyName string, validDurationInSeconds int64, bucketId string, namePrefix string) (*ApplicationKey, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_create_key", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId              string   `json:"accountId"`
			Capabilities           []string `json:"capabilities"`
			KeyName                string   `json:"keyName"`
			ValidDurationInSeconds int64    `json:"validDurationInSeconds,omitempty"`
			BucketId               string   `json:"bucketId,omitempty"`
			NamePrefix             string   `json:"namePrefix,omitempty"`
		}{b.auth.AccountId, capabilities, keyName, validDurationInSeconds, bucketId, namePrefix}
		responseBody = &ApplicationKey{}
	)

	response, err := b.makeAuthedRequest(url, requestBody)
	if err != nil {
		return nil, err
	}

	switch {
	case response.StatusCode == 200:
		if err = unmarshalResponseBody(response, responseBody); err != nil {
			return nil, err
		}
		return responseBody, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}

// DeleteKey deletes the application key specified.
// See "b2_delete_key" for an introduction:
// https://www.backblaze.com/b2/docs/b2_delete_key.html
//
// Parameter key is required.
// DeleteKey return nil if delete successd or error if something error happened.
func (b *B2) DeleteKey(key *ApplicationKey) error {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_delete_key", b.auth.ApiUrl)
		requestBody = &struct {
			Application string `json:"applicationKeyId"`
		}{key.ApplicationKeyId}
		responseBody = &ApplicationKey{}
	)

	response, err := b.makeAuthedRequest(url, requestBody)
	if err != nil {
		return err
	}

	switch {
	case response.StatusCode == 200:
		if err = unmarshalResponseBody(response, responseBody); err != nil {
			return err
		}
		return nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return handleErrorResponse(response)
	default:
		return handleUnknownResponse(response)
	}
}

// ListKeys list application keys associated with an account.
// See "b2_list_keys" for an introduction:
// https://www.backblaze.com/b2/docs/b2_list_keys.html
// All parameters are optional.
// ListKeys return an list of application pointer and an error.
func (b *B2) ListKeys(maxKeyCount int64, startApplicationKeyId string) (*ApplicationKeys, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_keys", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId             string `json:"accountId"`
			MaxKeyCount           int64  `json:"maxKeyCount,omitempty"`
			StartApplicationKeyId string `json:"startApplicationKeyId,omitempty"`
		}{b.auth.AccountId, maxKeyCount, startApplicationKeyId}
		responseBody = &ApplicationKeys{}
	)

	response, err := b.makeAuthedRequest(url, requestBody)
	if err != nil {
		return nil, err
	}

	switch {
	case response.StatusCode == 200:
		if err = unmarshalResponseBody(response, responseBody); err != nil {
			return nil, err
		}
		return responseBody, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}
