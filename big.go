// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// StartLargeFile register a file to upload
// See "b2_start_large_file" for an introduction:
// https://www.backblaze.com/b2/docs/b2_start_large_file.html
//
// Parameter bucketId and fileName are required. You can pass other empty value for other parameter for simplicity.
// StartLargeFile return a File array and an error.
func (b *B2) StartLargeFile(bucketId, fileName string, fileInfo map[string]string) (*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_start_large_file", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId    string            `json:"bucketId"`
			FileName    string            `json:"fileName"`
			ContentType string            `json:"contentType"`
			FileInfo    map[string]string `json:"fileInfo,omitempty"`
		}{bucketId, fileName, "b2/x-auto", fileInfo}
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

// GetUploadPartUrl return uploadUrlToken for a registered large file.
// See "b2_get_upload_part_url" for an introduction:
// https://www.backblaze.com/b2/docs/b2_get_upload_part_url.html
//
// Parameter fileId is required.
// GetUploadPartUrl return a File pointer and an error.
func (b *B2) GetUploadPartUrl(fileId string) (*UploadUrlToken, error) {
	var (
		url                     = fmt.Sprintf("%s/b2api/v1/b2_get_upload_part_url", b.auth.ApiUrl)
		getUploadPartUrlRequest = &struct {
			FileId string `json:"fileId"`
		}{FileId: fileId}
		uploadUrlToken = &UploadUrlToken{}
	)

	response, err := b.makeAuthedRequest(url, &getUploadPartUrlRequest)
	if err != nil {
		return nil, err
	}

	switch {
	case response.StatusCode == 200:
		if err := unmarshalResponseBody(response, uploadUrlToken); err != nil {
			return nil, err
		}
		return uploadUrlToken, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}
}

// UploadPart upload file part to b2 Close Storage.
// See "b2_upload_part" for an introduction:
// https://www.backblaze.com/b2/docs/b2_upload_part.html
//
// All parameters are required.
// UploadPart return a content sha1 and an error.
func (b *B2) UploadPart(uploadUrlToken *UploadUrlToken, filePath string, offset, size, partNumber int64,
	progress func(int64, int64)) (string, error) {
	contentSha1 := ""
	f, err := os.Open(filePath)
	if err != nil {
		return contentSha1, err
	}

	f.Seek(io.SeekStart, int(offset))
	buf := make([]byte, size)
	if _, err := f.Read(buf); err != nil {
		return contentSha1, err
	}

	headers := map[string]string{
		"X-Bz-Part-Number": strconv.FormatInt(partNumber, 10),
	}

	response, contentSha1, err := b.makeUploadRequest(uploadUrlToken, buf, headers, progress)
	if err != nil {
		return contentSha1, err
	}

	switch {
	case response.StatusCode == 200:
		return contentSha1, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return contentSha1, handleErrorResponse(response)
	default:
		return contentSha1, handleUnknownResponse(response)
	}
}

// ListParts return parts of a large file.
// See "b2_list_parts" for an introduction:
// https://www.backblaze.com/b2/docs/b2_list_parts.html
//
// All parameters are required.
// ListParts return uploaded parts of a large file.
func (b *B2) ListParts(fileId string, startPartNumber int64, maxPartNumber int64) ([]*Part, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_parts", b.auth.ApiUrl)
		requestBody = &struct {
			FileId          string `json:"fileId"`
			StartPartNumber int64  `json:"startPartNumber"`
			maxPartNumber   int64  `json:"maxPartNumber"`
		}{fileId, startPartNumber, maxPartNumber}
		responseBody = &struct {
			Parts          []*Part `json:"parts"`
			NextPartNumber int64   `json:"nextPartNumber"`
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
		return responseBody.Parts, nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return nil, handleErrorResponse(response)
	default:
		return nil, handleUnknownResponse(response)
	}

}

// ListUnfinishedLargeFiles lists unfinished large files.
// See "b2_list_unfinished_large_files" for an introduction:
// https://www.backblaze.com/b2/docs/b2_list_unfinished_large_files.html
//
// All parameters are required.
// ListUnfinishedLargeFiles return an array of file and an error.
func (b *B2) ListUnfinishedLargeFiles(bucketId, namePrefix string, startFileId string, maxfileCount int64) ([]*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_list_unfinished_large_files", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId     string `json:"bucketId"`
			NamePrefix   string `json:"namePrefix,omitempty"`
			StartFileId  string `json:"startFileId,omitempty"`
			MaxFileCount int64  `json:"maxFileCount,omitempty"`
		}{bucketId, namePrefix, startFileId, maxfileCount}
		responseBody = &struct {
			Files      []*File `json:"files"`
			NextFileId string  `json:"nextFileId"`
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

// FinishLargeFile merge all parts uploaded.
// See "b2_finish_large_file" for an introduction:
// https://www.backblaze.com/b2/docs/b2_finish_large_file.html
//
// Parameter fileId and partSha1Array are required.
// FinishLargeFile return a file pointer and an error.
func (b *B2) FinishLargeFile(fileId string, partSha1Array []string) (*File, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_finish_large_file", b.auth.ApiUrl)
		requestBody = &struct {
			FileId        string   `json:"fileId"`
			PartSha1Array []string `json:"partSha1Array"`
		}{fileId, partSha1Array}
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

// CancelLargeFile cancel parted uploaded.
// See "b2_cancel_large_file" for an introduction:
// https://www.backblaze.com/b2/docs/b2_cancel_large_file.html
//
// Parameter fileId is required.
// CancelLargeFile return an error.
func (b *B2) CancelLargeFile(fileId string) error {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_cancel_large_file", b.auth.ApiUrl)
		requestBody = &struct {
			FileId string `json:"fileId"`
		}{fileId}
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
