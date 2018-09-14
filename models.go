// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

type CorsRule struct {
	CorsRuleName      string   `json:"corsRuleName"`
	AllowedOrigins    []string `json:"allowedOrigins"`
	AllowedHeaders    []string `json:"allowedHeaders"`
	AllowedOperations []string `json:"allowedOperations"`
	ExposeHeaders     []string `json:"exposeHeaders"`
	MaxAgeSeconds     int64    `json:"maxAgeSeconds"`
}

type LifecycleRule struct {
	DaysFromHidingToDeleting  int64  `json:"daysFromHidingToDeleting,omitempty"`
	DaysFromUploadingToHiding int64  `json:"daysFromUploadingToHiding,omitempty"`
	FileNamePrefix            string `json:"fileNamePrefix"`
}

type Bucket struct {
	BucketId       string            `json:"bucketId"`
	BucketName     string            `json:"bucketName"`
	BucketType     string            `json:"bucketType"`
	BucketInfo     map[string]string `json:"bucketInfo,omitempty"`
	CorsRules      []CorsRule        `json:"corsRules,omitempty"`
	LifecycleRules []LifecycleRule   `json:"lifecycleRules,omitempty"`
	Revision       int64             `json:"revision,omitempty"`
}

const PUBLIC = "allPublic"
const PRIVATE = "allPrivate"

// Capability
const LIST_KEYS = "listKeys"
const WRITE_KEYS = "writeKeys"
const DELETE_KEYS = "deleteKeys"
const LIST_BUCKETS = "listBuckets"
const WRITE_BUCKETS = "writeBuckets"
const DELETE_BUCKETS = "deleteBuckets"
const LIST_FILES = "listFiles"
const READ_FILES = "readFiles"
const SHARE_FILES = "shareFiles"
const WRITE_FILES = "writeFiles"
const DELETE_FILES = "deleteFiles"

type ApplicationKey struct {
	ApplicationKeyId    string   `json:"applicationKeyId"`
	ApplicationKey      string   `json:"applicationKey"`
	KeyName             string   `json:"keyName"`
	Capabilities        []string `json:"capabilities"`
	ExpirationTimestamp int64    `json:"expirationTimestamp,omitempty"`
	BucketId            string   `json:"bucketId,omitempty"`
	NamePrefix          string   `json:"namePrefix,omitempty"`
}

type ApplicationKeys struct {
	Keys                 []*ApplicationKey `json:"keys"`
	NextApplicationKeyId string            `json:"nextApplicationKeyId,omitempty"`
}

type FileInfo map[string]interface{}

type File struct {
	FileId          string   `json:"fileId"`
	FileName        string   `json:"fileName"`
	ContentLength   int64    `json:"contentLength"`
	ContentType     string   `json:"contentType"`
	ContentSha1     string   `json:"contentSha1"`
	FileInfo        FileInfo `json:"fileInfo"`
	Action          string   `json:"action"`
	UploadTimestamp int64    `json:"uploadTimestamp"`
}

type DownloadUrlToken struct {
	FileNamePrefix     string `json:"fileNamePrefix"`
	AuthorizationToken string `json:"authorizationToken"`
}

type UploadUrlToken struct {
	UploadUrl          string `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

type Part struct {
	FileId          string `json:"fileId"`
	PartNumber      int64  `json:"partNumber"`
	ContentType     int64  `json:"contentType"`
	ContentSha1     string `json:"contentSha1"`
	UploadTimestamp int64  `json:"uploadTimestamp"`
}
