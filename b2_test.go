// Copyright 2018 hryyan. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package b2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	b2   *B2
	once sync.Once
)

type FileTests struct {
	Test       *testing.T
	BucketName string
}

const FILE = "testfile"

func (t *FileTests) RemoveFile(fileName string) {
	if err := os.Remove(fileName); err != nil {
		log.Println(err.Error())
		t.Test.Fatalf("Delete file %s failed", fileName)
	} else {
		log.Printf("Delete file %s successed!\n", fileName)
	}
}

func (t *FileTests) TestSmallFile() {
	// create bucket
	bucket, err := b2.CreateBucket(
		t.BucketName, PRIVATE,
		map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		[]CorsRule{{
			CorsRuleName: "downloadFromAnyOrigin",
			AllowedOrigins: []string{
				"https",
			},
			AllowedHeaders: []string{
				"range",
			},
			AllowedOperations: []string{
				"b2_download_file_by_name",
			},
			ExposeHeaders: []string{},
			MaxAgeSeconds: 3600,
		}},
		[]LifecycleRule{{
			DaysFromHidingToDeleting:  1,
			DaysFromUploadingToHiding: 10,
			FileNamePrefix:            "t",
		}})

	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Create bucket failed!")
	} else {
		log.Printf("Create bucket %s successed!\n", bucket.BucketName)
	}

	// delete bucket
	defer func() {
		if err = b2.DeleteBucket(bucket.BucketId); err != nil {
			t.Test.Fatal("Delete bucket failed!")
		} else {
			log.Printf("Deleted bucket %s.\n", bucket.BucketName)
		}
	}()

	// get upload url
	// b2.GetUploadUrl("")
	uploadUrlToken, err := b2.GetUploadUrl(bucket.BucketId)
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Get upload url failed!")
	} else {
		log.Printf("Get upload url successed, url: %s, token: %s!\n",
			uploadUrlToken.UploadUrl,
			uploadUrlToken.AuthorizationToken)
	}

	// create file version1
	f, err := os.Create(FILE)
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Create file(version1) failed!")
	} else {
		log.Println("Create file(version1) successed!")
	}
	fileV1Content := []byte("b2 test version1")
	f.Write(fileV1Content)
	f.Close()

	// delete file1 after test
	defer t.RemoveFile(FILE)

	var uploaded int64 = 0
	var fileSize int64 = 0
	go func() {
		var percent float64
		for {
			switch {
			case fileSize == 0:
			case uploaded != fileSize:
				percent = float64(uploaded) / float64(fileSize)
				log.Printf("Upload %.2f%%.\n", percent*100)
			case uploaded == fileSize:
				log.Println("Upload 100%.")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// upload file1 version1
	fileV1, err := b2.UploadFile(uploadUrlToken, FILE, func(done int64, total int64) {
		uploaded, fileSize = done, total
	})
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Upload file failed!")
	} else {
		log.Println("Upload file to bucket successed!")
	}

	defer func() {
		if err = b2.DeleteFileVersion(fileV1.FileName, fileV1.FileId); err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Delete file version1 failed!")
		} else {
			log.Println("Delete file version1 successed!")
		}
	}()

	// create file1 version2
	f, err = os.Create(FILE)
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Create file(version2) failed!")
	} else {
		log.Println("Create file(version2) successed!")
	}
	fileV2Content := []byte("b2 test version2")
	f.Write(fileV2Content)
	f.Close()

	uploaded, fileSize = 0, 0
	go func() {
		var percent float64
		for {
			switch {
			case fileSize == 0:
			case uploaded != fileSize:
				percent = float64(uploaded) / float64(fileSize)
				log.Printf("Upload %.2f%%.\n", percent*100)
			case uploaded == fileSize:
				log.Println("Upload 100%.")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// upload file1 version2
	fileV2, err := b2.UploadFile(uploadUrlToken, FILE, func(done int64, total int64) {
		uploaded, fileSize = done, total
	})
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Upload file failed!")
	} else {
		log.Println("Upload file to bucket successed!")
	}

	defer func() {
		if err = b2.DeleteFileVersion(fileV2.FileName, fileV2.FileId); err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Delete file version2 failed!")
		} else {
			log.Println("Delete file version2 successed!")
		}
	}()

	// list file versions
	if _, err := b2.ListFileVersions(bucket.BucketId, "", "", "", "", 1000); err != nil {
		log.Println(err.Error())
		t.Test.Fatal("List file versions failed!")
	} else {
		log.Println("List file versions successed!")
	}

	// list file names
	if _, err = b2.ListFileNames(bucket.BucketId, "", "", "", 1000); err != nil {
		log.Println(err.Error())
		t.Test.Fatal("List file names failed!")
	} else {
		log.Println("List file names successed!")
	}

	// get file version1 info
	fileV1Info, err := b2.GetFileInfo(fileV1.FileId)
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Get file version1 info failed!")
	} else {
		log.Println("Get file version1 info successed!")
	}

	// get file version2 info
	fileV2Info, err := b2.GetFileInfo(fileV2.FileId)
	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Get file version2 info failed!")
	} else {
		log.Println("Get file version2 info successed!")
	}

	// log.Println("Sleep 10 seconds for correct version!")
	// time.Sleep(10 * time.Second)

	// download file by id
	var downloaded int64
	fileSize = 0
	go func() {
		var percent float64
		for {
			switch {
			case fileSize == 0:
			case downloaded != fileSize:
				percent = float64(downloaded) / float64(fileSize)
				log.Printf("Download %.2f%%.\n", percent*100)
			case downloaded == fileSize:
				log.Println("Download 100%.")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	fileName := fmt.Sprintf("%s.v1", FILE)
	if err = b2.DownloadFileById(fileV1Info.FileId, fileName,
		true, func(done int64, total int64) {
			downloaded, fileSize = done, total
		}); err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Download file version1 failed(by file version)!")
	} else {
		log.Println("Download file version1 successed(by file version)!")
		fileV1, err := os.Open(fileName)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Cannot open file version1")
		} else {
			b, err := ioutil.ReadAll(fileV1)
			if err != nil {
				log.Println(err.Error())
				t.Test.Fatal("Cannot read file version1")
			}
			if !bytes.Equal(b, fileV1Content) {
				t.Test.Fatal("File verion1 content dismatch!")
			} else {
				log.Println("File version1 content corrected!")
			}
		}
		defer t.RemoveFile(fileName)
	}

	// download file by name
	downloaded, fileSize = 0, 0
	go func() {
		var percent float64
		for {
			switch {
			case fileSize == 0:
			case downloaded != fileSize:
				percent = float64(downloaded) / float64(fileSize)
				log.Printf("Download %.2f%%.\n", percent*100)
			case downloaded == fileSize:
				log.Println("Download 100%.")
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	fileName = fmt.Sprintf("%s.v2", FILE)
	if err = b2.DownloadFileByName(bucket.BucketName, fileV2Info.FileName,
		fileName, true, func(done int64, total int64) {
			downloaded, fileSize = done, total
		}); err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Download file version2 failed(by file name)!")
	} else {
		log.Println("Download file version1 successed(by file name)!")
		_, err := os.Open(fileName)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Cannot open file version2")
		} else {
			// b, err := ioutil.ReadAll(fileV2)
			// if err != nil {
			// 	log.Println(err.Error())
			// 	t.Test.Fatal("Cannot read file version2")
			// }
			// if !bytes.Equal(b, fileV2Content) {
			// 	t.Test.Fatal("File verion2 content dismatch!")
			// } else {
			// 	log.Println("File version2 content corrected!")
			// }
		}
		defer t.RemoveFile(fileName)
	}

	// // hide file
	// if err = b2.HideFile(bucket.BucketId, fileV1Info.FileName); err != nil {
	// 	log.Println(err.Error())
	// 	t.Test.Fatal("Hide file version1 failed!")
	// } else {
	// 	log.Println("Hide file version1 successed!")
	// }
}

func (t *FileTests) TestLargeFile() {
	// create bucket
	bucket, err := b2.CreateBucket(
		t.BucketName, PRIVATE,
		map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		[]CorsRule{{
			CorsRuleName: "downloadFromAnyOrigin",
			AllowedOrigins: []string{
				"https",
			},
			AllowedHeaders: []string{
				"range",
			},
			AllowedOperations: []string{
				"b2_download_file_by_name",
			},
			ExposeHeaders: []string{},
			MaxAgeSeconds: 3600,
		}},
		[]LifecycleRule{{
			DaysFromHidingToDeleting:  1,
			DaysFromUploadingToHiding: 10,
			FileNamePrefix:            "t",
		}})

	if err != nil {
		log.Println(err.Error())
		t.Test.Fatal("Create bucket failed!")
	} else {
		log.Printf("Create bucket %s successed!\n", bucket.BucketName)
	}

	// delete bucket
	defer func() {
		if err = b2.DeleteBucket(bucket.BucketId); err != nil {
			t.Test.Fatal("Delete bucket failed!")
		} else {
			log.Printf("Deleted bucket %s.\n", bucket.BucketName)
		}
	}()

	{
		// start large file
		file, err := b2.StartLargeFile(bucket.BucketId, FILE, map[string]string{})
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Start large file failed!")
		} else {
			log.Println("Start large file successed!")
		}

		// get upload part url
		uploadUrlToken, err := b2.GetUploadPartUrl(file.FileId)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Get upload part url failed!")
		} else {
			log.Println("Get upload part url successed!")
		}

		// create file version1
		f, err := os.Create(FILE)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Create file failed!")
		} else {
			log.Println("Create file successed!")
		}
		if err = f.Truncate(10000000); err != nil {
			log.Println("Truncate file failed!")
		}
		f.Close()
		defer t.RemoveFile(FILE)

		// upload part
		var uploaded int64 = 0
		var fileSize int64 = 0
		go func() {
			var percent float64
			for {
				switch {
				case fileSize == 0:
				case uploaded != fileSize:
					percent = float64(uploaded) / float64(fileSize)
					log.Printf("Upload %.2f%%.\n", percent*100)
				case uploaded == fileSize:
					log.Println("Upload 100%.")
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		part1ContentSha1, err := b2.UploadPart(uploadUrlToken, FILE, 0, 5000000, 1, func(done int64, total int64) {
			uploaded, fileSize = done, total
		})
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Upload part 1 failed!")
		} else {
			log.Println("Upload part 1 successed!")
		}

		uploaded, fileSize = 0, 0
		go func() {
			var percent float64
			for {
				switch {
				case fileSize == 0:
				case uploaded != fileSize:
					percent = float64(uploaded) / float64(fileSize)
					log.Printf("Upload %.2f%%.\n", percent*100)
				case uploaded == fileSize:
					log.Println("Upload 100%.")
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		part2ContentSha1, err := b2.UploadPart(uploadUrlToken, FILE, 5000000, 10000000, 2, func(done int64, total int64) {
			uploaded, fileSize = done, total
		})
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Upload part 2 failed!")
		} else {
			log.Println("Upload part 2 successed!")
		}

		// list part
		_, err = b2.ListParts(file.FileId, 1, 1000)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("List parts failed!")
		} else {
			log.Println("List parted successed!")
		}

		// finish large file
		partSha1Array := []string{part1ContentSha1, part2ContentSha1}
		file, err = b2.FinishLargeFile(file.FileId, partSha1Array)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Finish parts failed!")
		} else {
			log.Println("Finish parts successed!")
		}

		// download large file
		var downloaded int64 = 0
		fileSize = 0
		go func() {
			var percent float64
			for {
				switch {
				case fileSize == 0:
				case downloaded != fileSize:
					percent = float64(downloaded) / float64(fileSize)
					log.Printf("Download %.2f%%.\n", percent*100)
				case downloaded == fileSize:
					log.Println("Download 100%.")
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()

		fileName := fmt.Sprintf("%s.golden", FILE)
		defer t.RemoveFile(fileName)

		if err = b2.DownloadFileByName(bucket.BucketName, file.FileName, fileName,
			true, func(done int64, total int64) {
				downloaded, fileSize = done, total
			}); err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Download large file failed(by file name!)")
		}

		// delete file version
		if err = b2.DeleteFileVersion(file.FileName, file.FileId); err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Delete file version failed!")
		} else {
			log.Println("Delete file version successed!")
		}
	}

	{
		// start large file
		FILE := FILE + "2"
		file, err := b2.StartLargeFile(bucket.BucketId, FILE, map[string]string{})
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Start large file failed!")
		} else {
			log.Println("Start large file successed!")
		}

		// get upload part url
		uploadUrlToken, err := b2.GetUploadPartUrl(file.FileId)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Get upload part url failed!")
		} else {
			log.Println("Get upload part url successed!")
		}

		// create file version1
		f, err := os.Create(FILE)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Create file failed!")
		} else {
			log.Println("Create file successed!")
		}
		if err = f.Truncate(5000000); err != nil {
			log.Println("Truncate file failed!")
		}
		f.Close()
		defer t.RemoveFile(FILE)

		// upload part
		var uploaded int64 = 0
		var fileSize int64 = 0
		go func() {
			var percent float64
			for {
				switch {
				case fileSize == 0:
				case uploaded != fileSize:
					percent = float64(uploaded) / float64(fileSize)
					log.Printf("Upload %.2f%%.\n", percent*100)
				case uploaded == fileSize:
					log.Println("Upload 100%.")
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()
		_, err = b2.UploadPart(uploadUrlToken, FILE, 0, 5000000,
			1, func(done int64, total int64) {
				uploaded, fileSize = done, total
			})
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Upload part 1 failed!")
		} else {
			log.Println("Upload part 1 successed!")
		}

		_, err = b2.ListParts(file.FileId, 1, 1000)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("List parts failed!")
		} else {
			log.Println("List parts successed!")
		}

		// list unfinished large file
		_, err = b2.ListUnfinishedLargeFiles(bucket.BucketId, "", "", 100)
		if err != nil {
			log.Println(err.Error())
			t.Test.Fatal("List unfinished large files failed!")
		} else {
			log.Println("List unfinished large files successed!")
		}

		// cancel large file
		if err = b2.CancelLargeFile(file.FileId); err != nil {
			log.Println(err.Error())
			t.Test.Fatal("Cancel large file failed")
		} else {
			log.Println("Cancel large file successed!")
		}
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randName(length int64) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestFile(t *testing.T) {
	t.Run("A=create", func(t *testing.T) {
		test := FileTests{Test: t, BucketName: randName(16)}
		test.TestSmallFile()
		test.TestLargeFile()
	})
}

func getKeyFromEnv() (string, string) {
	return os.Getenv("B2_ACCOUNT_ID"), os.Getenv("B2_APPLICATION_KEY")
}

func setup() {
	once.Do(func() {
		accountId, applicationKey := getKeyFromEnv()
		b2 = &B2{
			AccountId:      accountId,
			ApplicationKey: applicationKey,
		}

		if err := b2.Auth(); err != nil {
			log.Fatal("Authorization failed!")
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	runTests := m.Run()
	os.Exit(runTests)
}
