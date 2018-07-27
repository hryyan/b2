// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetUploadUrl hide the file.
// See "b2_get_upload_url" for an introduction:
// https://www.backblaze.com/b2/docs/b2_get_upload_url.html
//
// Parameter bucketId is required, you can get bucketId from any Bucket struct.
// GetUploadUrl return a UploadUrlToken pointer and an error.
func (b *B2) GetUploadUrl(bucketId string) (*UploadUrlToken, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_get_upload_url", b.auth.ApiUrl)
		requestBody = struct {
			BucketId string `json:"bucketId"`
		}{BucketId: bucketId}
		responseBody = &UploadUrlToken{}
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

// UploadFile upload file to b2 Close Storage.
// See "b2_upload_file" for an introduction:
// https://www.backblaze.com/b2/docs/b2_upload_file.html
//
// Parameter uploadUrlToken and filePath are required.
// UploadFile return a File pointer and an error.
func (b *B2) UploadFile(uploadUrlToken *UploadUrlToken, filePath string, progress func(int64, int64)) (*File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fi.Size())
	if _, err := f.Read(buf); err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-Bz-File-Name":                     filepath.Base(filePath),
		"X-Bz-Info-src_last_modified_millis": fmt.Sprintf("%d", fi.ModTime().Unix()*1000),
	}

	response, _, err := b.makeUploadRequest(uploadUrlToken, buf, headers, progress)
	if err != nil {
		return nil, err
	}

	switch {
	case response.StatusCode == 200:
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		var file File
		if err = json.Unmarshal(body, &file); err != nil {
			return nil, err
		}

		return &file, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}
