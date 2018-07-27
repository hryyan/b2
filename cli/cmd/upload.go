package cmd

import (
	"fmt"
	"math"
	"os"
	"path"
	"sync"

	"github.com/hryyan/b2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var concurrency int64 = 1

func uploadFile(client *b2.B2, bucket *b2.Bucket, fileName, filePath string, size int64) {
	var (
		wg sync.WaitGroup
		p  = mpb.New(mpb.WithWaitGroup(&wg))
	)

	uploadUrlToken, err := client.GetUploadUrl(bucket.BucketId)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(B2_LIBRARY_ERROR_EXIT)
	}

	bar := p.AddBar(
		size,
		mpb.PrependDecorators(
			decor.Name(fileName, decor.WC{W: len(fileName), C: decor.DidentRight}),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.AverageSpeed(decor.UnitKiB, "% .2f "),
			decor.Name("eta: ", decor.WC{W: 5}),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO), "Done",
			),
		),
	)

	go func() {
		wg.Add(1)
		defer wg.Done()

		var offset int64
		if _, err = client.UploadFile(uploadUrlToken, filePath, func(done, total int64) {
			bar.IncrBy(int(done - offset))
			offset = done
		}); err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}
	}()

	p.Wait()
}

func uploadParts(client *b2.B2, bucket *b2.Bucket, fileName, filePath string, size, concurrency int64) {
	var (
		start     int64 = 0
		partSize  int64 = size / concurrency
		sha1Array       = make([]string, concurrency)
		i         int64 = 0

		wg sync.WaitGroup
		p  = mpb.New(mpb.WithWaitGroup(&wg))
	)

	file, err := client.StartLargeFile(bucket.BucketId, fileName, map[string]string{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(B2_LIBRARY_ERROR_EXIT)
	}

	for i = 0; i < concurrency; i++ {
		uploadUrlToken, err := client.GetUploadPartUrl(file.FileId)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		if i == concurrency-1 {
			partSize = size - partSize*i
		}

		bar := p.AddBar(
			partSize,
			mpb.PrependDecorators(
				decor.Name(fileName, decor.WC{W: len(fileName), C: decor.DidentRight}),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.AverageSpeed(decor.UnitKiB, "% .2f "),
				decor.Name("eta: ", decor.WC{W: 5}),
				decor.OnComplete(
					decor.AverageETA(decor.ET_STYLE_GO), "Done",
				),
			),
		)

		go func(start, pageSize, index int64) {
			wg.Add(1)
			defer wg.Done()

			var offset int64
			contentSha1, err := client.UploadPart(uploadUrlToken, filePath, start, partSize, index+1, func(done, total int64) {
				bar.IncrBy(int(done - offset))
				offset = done
			})

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(B2_LIBRARY_ERROR_EXIT)
			}

			sha1Array[index] = contentSha1
		}(start, partSize, i)

		start += partSize
	}

	p.Wait()

	if _, err := client.FinishLargeFile(file.FileId, sha1Array); err != nil {
		fmt.Println(err.Error())
		os.Exit(B2_LIBRARY_ERROR_EXIT)
	}
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload bucket file",
	Short: "Upload file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			bucketName = args[0]
			fileName   = args[1]
			filePath   = ""
		)

		if path.IsAbs(fileName) {
			filePath = fileName
			fileName = path.Base(filePath)
		} else {
			filePath = path.Join(".", fileName)
		}

		info, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Read file info error!")
			os.Exit(OPERATION_ERROR_EXIT)
		}
		size := info.Size()

		if concurrency < 1 {
			concurrency = 1
		}

		minPart, maxPart := 5000000.0, 500000000.0
		part := float64(size / concurrency)
		suggestConcurrency := int64(math.Ceil(float64(size) / 100000000.0))

		if part < minPart {
			fmt.Println("Below min part size(5M), auto set part size to 100M!")
			fmt.Println("Set concurrency to ", suggestConcurrency)
			concurrency = suggestConcurrency
		} else if part > maxPart {
			fmt.Println("Above max part size(500M), auto set part size to 100M!")
			fmt.Println("Set concurrency to ", suggestConcurrency)
			concurrency = suggestConcurrency
		}

		client := login()
		buckets, err := client.ListBuckets("", bucketName, "")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		if len(buckets) != 1 {
			fmt.Println("Can not find bucket %s!\n", bucketName)
			os.Exit(OPERATION_ERROR_EXIT)
		}

		bucket := buckets[0]

		if concurrency == 1 {
			uploadFile(client, bucket, fileName, filePath, size)
		} else {
			uploadParts(client, bucket, fileName, filePath, size, concurrency)
		}
	},
}

func init() {
	uploadFileCmd.Flags().Int64VarP(
		&concurrency,
		"concurrency",
		"c",
		1,
		"threads for uploading")

	rootCmd.AddCommand(uploadFileCmd)
}
