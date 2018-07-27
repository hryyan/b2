// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"fmt"
)

// CreateBucket create a bucket.
// See "b2_create_bucket" for an introduction:
// https://www.backblaze.com/b2/docs/b2_create_bucket.html
//
// Parameter bucketName and bucketType are required. You can pass empty map or slice for simplicity.
// CreateBucket returned a Bucket pointer and an error.
func (b *B2) CreateBucket(bucketName, bucketType string, bucketInfo map[string]string,
	corsRules []CorsRule, lifecycleRules []LifecycleRule) (*Bucket, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_create_bucket", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId      string            `json:"accountId"`
			BucketName     string            `json:"bucketName"`
			BucketType     string            `json:"bucketType"`
			BucketInfo     map[string]string `json:"bucketInfo,omitempty"`
			CorsRules      []CorsRule        `json:"corsRules,omitempty"`
			LifecycleRules []LifecycleRule   `json:"lifecycleRules,omitempty"`
		}{b.AccountId, bucketName, bucketType, bucketInfo, corsRules, lifecycleRules}
		responseBody = &Bucket{}
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

// DeleteBucket delete a bucket.
// See "b2_delete_bucket" for an introduction:
// https://www.backblaze.com/b2/docs/b2_delete_bucket.html
//
// Parameter bucketId is required, you can get bucketId from an Bucket struct.
// DeleteBucket returned nil if success, return error if failed.
func (b *B2) DeleteBucket(bucketId string) error {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_delete_bucket", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId string `json:"accountId"`
			BucketId  string `json:"bucketId,omitempty"`
		}{b.AccountId, bucketId}
	)

	response, err := b.makeAuthedRequest(url, requestBody)
	if err != nil {
		return err
	}

	switch {
	case response.StatusCode == 200:
		return nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return handleErrorResponse(response)
	default:
		return handleUnknownResponse(response)
	}

}

// UpdateBucket update a bucket.
// See "b2_update_bucket" for an introduction:
// https://www.backblaze.com/b2/docs/b2_update_bucket.html
//
// Parameter bucket is required, you can pass ifRevisionIs as false for simplicity.
// UpdateBucket returned a bucket pointer and an error.
func (b *B2) UpdateBucket(bucket *Bucket, ifRevisionIs bool) (*Bucket, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_update_bucket", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId      string            `json:"accountId"`
			BucketId       string            `json:"bucketId"`
			BucketInfo     map[string]string `json:"bucketInfo,omitempty"`
			BucketType     string            `json:"bucketType,omitempty"`
			CorsRules      []CorsRule        `json:"corsRules,omitempty"`
			LifecycleRules []LifecycleRule   `json:"lifecycleRules,omitempty"`
			IfRevisionIs   bool              `json:"ifRevisionIs,omitempty"`
		}{
			b.AccountId,
			bucket.BucketId,
			bucket.BucketInfo,
			bucket.BucketType,
			bucket.CorsRules,
			bucket.LifecycleRules,
			ifRevisionIs,
		}
		responseBody = &Bucket{}
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

// ListBuckets update a bucket.
// See "b2_list_buckets" for introduction:
// https://www.backblaze.com/b2/docs/b2_list_buckets.html
//
// All parameter are optional, you can pass empty value for simplicity.
// List returned a bucket array and an error.
func (b *B2) ListBuckets(bucketId, bucketName, bucketTypes string) ([]*Bucket, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_buckets", b.auth.ApiUrl)
		requestBody = &struct {
			AccountId   string `json:"accountId"`
			BucketId    string `json:"bucketId,omitempty"`
			BucketName  string `json:"bucketName,omitempty"`
			BucketTypes string `json:"bucketTypes,omitempty"`
		}{b.AccountId, bucketId, bucketName, bucketTypes}
		responseBody = &struct {
			Buckets []*Bucket `json: "buckets"`
		}{}
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
		return responseBody.Buckets, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}
