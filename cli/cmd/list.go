package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hryyan/b2"
	"github.com/spf13/cobra"
)

var listBucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "List buckets",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := login()
		buckets, err := client.ListBuckets("", "", "")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}
		green := color.New(color.FgGreen).PrintfFunc()
		white := color.New(color.FgWhite).PrintfFunc()

		for _, bucket := range buckets {
			if bucket.BucketType == b2.PUBLIC {
				green("%s\n", bucket.BucketName)
			} else {
				white("%s\n", bucket.BucketName)
			}
		}
	},
}

var listFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "List file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := login()
		bucketName := args[0]
		buckets, err := client.ListBuckets("", bucketName, "")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		if len(buckets) != 1 {
			fmt.Printf("Can not find bucket %s!\n", bucketName)
			os.Exit(OPERATION_ERROR_EXIT)
		}

		bucket := buckets[0]
		files, err := client.ListFileNames(bucket.BucketId, "", "", "", 10000)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		for _, file := range files {
			fmt.Printf("%s %d %s\n", file.FileName, file.ContentLength, file.ContentType)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list command",
	Short: "List buckets or files",
}

func init() {
	listCmd.AddCommand(listBucketsCmd, listFilesCmd)

	rootCmd.AddCommand(listCmd)
}
