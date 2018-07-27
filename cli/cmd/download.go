package cmd

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/hryyan/b2"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var saveTo string

const MAX_VALID_DURATION_IN_SECONDS = 604800

func downloadFile(client *b2.B2, bucket *b2.Bucket, fileName, filePath string) {
	var (
		wg sync.WaitGroup
		p  = mpb.New(mpb.WithWaitGroup(&wg))
	)

	bar := p.AddBar(
		0,
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
		if err := client.DownloadFileByName(bucket.BucketName, fileName, filePath, true, func(done, total int64) {
			bar.SetTotal(total, false)
			bar.IncrBy(int(done - offset))
			offset = done
		}); err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}
	}()

	p.Wait()
}

var downloadFileCmd = &cobra.Command{
	Use:   "download bucket file",
	Short: "Download file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			bucketName = args[0]
			fileName   = args[1]
			filePath   = ""
		)
		if saveTo != "" {
			filePath = saveTo
		} else {
			filePath = path.Join(".", fileName)
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

		downloadFile(client, bucket, fileName, filePath)
	},
}

func init() {
	downloadFileCmd.Flags().StringVarP(
		&saveTo,
		"save",
		"s",
		"",
		"save as file")

	rootCmd.AddCommand(downloadFileCmd)
}
