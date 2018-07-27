// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"fmt"
	"log"
)

// ListFileNames list files names.
// See "b2_list_file_names" for an introduction:
// https://www.backblaze.com/b2/docs/b2_list_file_names.html
//
// Parameter bucketId is required. You can pass other empty value for other parameter for simplicity.
// ListFileNames return a File array and an error.
func (b *B2) ListFileNames(bucketId, startFileName, prefix, delimiter string,
	maxFileCount int64) ([]*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_file_names", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId      string `json:"bucketId"`
			StartFileName string `json:"startFileName,omitempty"`
			Prefix        string `json:"prefix,omitempty"`
			Delimiter     string `json:"delimiter,omitempty"`
			MaxFileCount  int64  `json:"maxFileCount,omitempty"`
		}{bucketId, startFileName, prefix, delimiter, maxFileCount}
		responseBody = &struct {
			Files        []*File `json:"files"`
			NextFileName string  `json:"nextFileName"`
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
		return responseBody.Files, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}

// ListFileVersions list files versions.
// See "b2_list_file_versions" for an introduction:
// https://www.backblaze.com/b2/docs/b2_list_file_versions.html
//
// Parameter bucketId is required. You can pass other empty value for other parameter for simplicity.
// ListFileVersions return a File array and an error.
func (b *B2) ListFileVersions(bucketId, startFileName, startFileId, prefix, delimiter string,
	maxFileCount int64) ([]*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_file_versions", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId      string `json:"bucketId"`
			StartFileName string `json:"startFileName,omitempty"`
			StartFileId   string `json:"startField,omitempty"`
			Prefix        string `json:"prefix,omitempty"`
			Delimiter     string `json:"delimiter,omitempty"`
			MaxFileCount  int64  `json:"maxFileCount,omitempty"`
		}{bucketId, startFileName, startFileId, prefix, delimiter, maxFileCount}
		responseBody = &struct {
			Files        []*File `json:"files"`
			NextFileName string  `json:"nextFileName"`
			NextField    string  `json:"nextField"`
		}{}
	)

	response, err := b.makeAuthedRequest(url, requestBody)
	if err != nil {
		return nil, err
	}

	log.Println(url)

	switch {
	case response.StatusCode == 200:
		if err = unmarshalResponseBody(response, responseBody); err != nil {
			return nil, err
		}
		return responseBody.Files, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}

// GetFileInfo get the file info.
// See "b2_get_file_info" for an introduction:
// https://www.backblaze.com/b2/docs/b2_get_file_info.html
//
// Parameter fileId is required, you can get fileId from any File struct.
// GetFileInfo return a File pointer and an error.
func (b *B2) GetFileInfo(fileId string) (*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_get_file_info", b.auth.ApiUrl)
		requestBody = &struct {
			FileId string `json:"fileId"`
		}{fileId}
		responseBody = &File{}
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

// HideFile hide the file.
// See "b2_hide_file" for an introduction:
// https://www.backblaze.com/b2/docs/b2_hide_file.html
//
// Parameter bucketId and fileName is required.
// HideFile return nil if successed, return error if failed.
func (b *B2) HideFile(bucketId, fileName string) error {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_hide_file", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId string `json:"bucketId"`
			FileName string `json:"fileName"`
		}{bucketId, fileName}
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

// DeleteFileVersion delete a file version from b2 Cloud Storage. If the version you delete is the lastest version, and there are older versions, then the most recent older version will become the current version.
// See "b2_delete_file_version" for an introduction:
// https://www.backblaze.com/b2/docs/b2_delete_file_version.html
//
// Parameter fileName and fileId are required, you can get these from any File struct.
// DeleteFileVersion return nil if successed, return error if failed.
func (b *B2) DeleteFileVersion(fileName, fileId string) error {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_delete_file_version", b.auth.ApiUrl)
		requestBody = &struct {
			FileName string `json:"fileName"`
			FileId   string `json:"fileId"`
		}{fileName, fileId}
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
