package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hryyan/b2"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share bucket file",
	Short: "Share file and generate download url",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			client     = login()
			bucketName = args[0]
			fileName   = args[1]
		)

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
		if bucket.BucketType != b2.PUBLIC {
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Printf(`Bucket %s is not public, `+
					`do you want to publicise it? Type y / n`+
					"\n",
					bucket.BucketName)
				text, _ := reader.ReadString('\n')
				text = strings.Trim(text, "\n")
				if text == "y" {
					break
				} else if text == "n" {
					return
				} else {
					fmt.Println("Please type y / n")
					continue
				}
			}

			bucket.BucketType = b2.PUBLIC
			if _, err = client.UpdateBucket(bucket, false); err != nil {
				fmt.Println(err.Error())
				os.Exit(B2_LIBRARY_ERROR_EXIT)
			}
		}

		URL := client.GetPublicFileDownloadURL(bucket.BucketName, fileName)
		fmt.Println(URL)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
