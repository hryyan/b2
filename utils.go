// Copyright 2018 hyyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type ProgressReaderWriter struct {
	Reader io.Reader
	Writer io.Writer
	Total  int64
	Done   int64
	Report func(int64, int64)
}

func (prw *ProgressReaderWriter) Read(p []byte) (int, error) {
	n, err := prw.Reader.Read(p)
	prw.Done += int64(n)
	prw.Report(prw.Done, prw.Total)
	return n, err
}

func (prw *ProgressReaderWriter) Write(p []byte) (int, error) {
	n, err := prw.Writer.Write(p)
	prw.Done += int64(n)
	prw.Report(prw.Done, prw.Total)
	return n, err
}

func (b *B2) makeAuthedRequest(url string, s interface{}) (*http.Response, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", b.auth.AuthorizationToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *B2) makeDownloadRequest(url string, queries map[string]string, needAuth bool) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	for key, value := range queries {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	if needAuth {
		request.Header.Set("Authorization", b.auth.AuthorizationToken)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *B2) makeUploadRequest(uploadUrlToken *UploadUrlToken, body []byte,
	headers map[string]string, report func(int64, int64)) (*http.Response, string, error) {
	var (
		h              = sha1.New()
		contentSha1    = ""
		progressReader = ProgressReaderWriter{
			Reader: bytes.NewReader(body),
			Total:  int64(len(body)),
			Report: report,
		}
	)

	if _, err := io.Copy(h, bytes.NewReader(body)); err != nil {
		return nil, contentSha1, err
	}

	request, err := http.NewRequest("POST", uploadUrlToken.UploadUrl, &progressReader)
	if err != nil {
		return nil, contentSha1, err
	}

	contentSha1 = fmt.Sprintf("%x", h.Sum(nil))

	request.ContentLength = int64(len(body))
	headers["Authorization"] = uploadUrlToken.AuthorizationToken
	headers["Content-Type"] = "b2/x-auto"
	headers["X-Bz-Content-Sha1"] = contentSha1
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, contentSha1, err
	}

	return response, contentSha1, nil
}

func unmarshalResponseBody(response *http.Response, s interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.Unmarshal(body, s); err != nil {
		return err
	}
	return nil
}

func GetKeyFromEnv() (string, string) {
	return os.Getenv("B2_ACCOUNT_ID"), os.Getenv("B2_APPLICATION_KEY")
}
