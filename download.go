// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"fmt"
	"io"
	"os"
)

// GetDownloadAuthorization create a download url and a token.
// See "b2_get_download_authorization" for an introduction:
// https://www.backblaze.com/b2/docs/b2_get_download_authorization.html
//
// Parameter bucketId, fileNamePrefix and validDurationInSeconds are required.
// GetDownloadAuthorization return a DownloadUrlToken pointer and an error.
func (b *B2) GetDownloadAuthorization(bucketId, fileNamePrefix string,
	validDurationInSeconds int64) (*DownloadUrlToken, error) {
	var (
		url         = fmt.Sprintf("%s/b2api/v1/b2_get_download_authorization", b.auth.ApiUrl)
		requestBody = &struct {
			BucketId               string `json:"bucketId"`
			FileNamePrefix         string `json:"fileNamePrefix"`
			ValidDurationInSeconds int64  `json:"validDurationInSeconds"`
		}{bucketId, fileNamePrefix, validDurationInSeconds}
		responseBody = &DownloadUrlToken{}
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

// DownloadFileById downlaod file from b2 Cloud Storage using fileId.
// See "b2_download_file_by_id" for an introduction:
// https://www.backblaze.com/b2/docs/b2_download_file_by_id.html
//
// Parameter fileId and filePath are required, if the bucket is private, you should pass needAuth as true.
// Parameter filePath is the local file path you want to save.
// DownloadFileById return nil if successed, return error if failed.
func (b *B2) DownloadFileById(fileId, filePath string, needAuth bool,
	report func(int64, int64)) error {
	var (
		url     = fmt.Sprintf("%s/b2api/v1/b2_download_file_by_id", b.auth.DownloadUrl)
		queries = map[string]string{
			"fileId": fileId,
		}
	)

	response, err := b.makeDownloadRequest(url, queries, needAuth)
	if err != nil {
		return err
	}

	switch {
	case response.StatusCode == 200:
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		defer response.Body.Close()

		progressWriter := &ProgressReaderWriter{
			Writer: f,
			Total:  response.ContentLength,
			Report: report,
		}

		if _, err = io.Copy(progressWriter, response.Body); err != nil {
			return err
		}
		return nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return handleErrorResponse(response)
	default:
		return handleUnknownResponse(response)
	}
}

// DownloadFileByName downlaod file from b2 Cloud Storage using bucketName and fileName.
// See "b2_download_file_by_name" for an introduction:
// https://www.backblaze.com/b2/docs/b2_download_file_by_name.html
//
// Parameter bucketName, fileName and filePath are required, if the bucket is private, you should pass needAuth as true.
// Parameter filePath is the local file path you want to save.
// DownloadFileByName return nil if successed, return error if failed.
func (b *B2) DownloadFileByName(bucketName, fileName, filePath string,
	needAuth bool, report func(int64, int64)) error {
	var (
		url = fmt.Sprintf("%s/file/%s/%s", b.auth.DownloadUrl, bucketName, fileName)
	)

	response, err := b.makeDownloadRequest(url, map[string]string{}, needAuth)
	if err != nil {
		return err
	}

	switch {
	case response.StatusCode == 200:
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		defer response.Body.Close()

		progressWriter := &ProgressReaderWriter{
			Writer: f,
			Total:  response.ContentLength,
			Report: report,
		}

		_, err = io.Copy(progressWriter, response.Body)
		if err != nil {
			return err
		}
		return nil
	case response.StatusCode == 400 || response.StatusCode == 401:
		return handleErrorResponse(response)
	default:
		return handleUnknownResponse(response)
	}
}

func (b *B2) GetPublicFileDownloadURL(bucketName, fileName string) string {
	return fmt.Sprintf("%s/file/%s/%s", b.auth.DownloadUrl, bucketName, fileName)
}
