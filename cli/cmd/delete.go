package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteBucketCmd = &cobra.Command{
	Use:   "bucket [bucket ..]",
	Short: "Delete bucket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := login()
		buckets, err := client.ListBuckets("", "", "")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		for _, name := range args {
			found := false
			for _, bucket := range buckets {
				if name == bucket.BucketName {
					found = true
					err := client.DeleteBucket(bucket.BucketId)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(OPERATION_ERROR_EXIT)
					}
				}
			}

			if !found {
				fmt.Printf("Can not find bucket %s!\n", name)
				os.Exit(OPERATION_ERROR_EXIT)
			}
		}
	},
}

var deleteFileCmd = &cobra.Command{
	Use:   "file bucket [file ..]",
	Short: "Delete file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := login()
		bucketName := args[0]
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
		files, err := client.ListFileNames(bucket.BucketId, "", "", "", 10000)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(B2_LIBRARY_ERROR_EXIT)
		}

		for _, arg := range args[1:] {
			found := false
			for _, file := range files {
				if file.FileName == arg {
					found = true
					if err = client.DeleteFileVersion(file.FileName, file.FileId); err != nil {
						fmt.Println(err.Error())
						os.Exit(B2_LIBRARY_ERROR_EXIT)
					} else {
						return
					}
				}
			}

			if !found {
				fmt.Printf("Can not find %s in %s!\n", args[1], args[0])
				os.Exit(OPERATION_ERROR_EXIT)
			}
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete command",
	Short: "Delete buckets or files",
}

func init() {
	deleteCmd.AddCommand(deleteBucketCmd, deleteFileCmd)

	rootCmd.AddCommand(deleteCmd)
}
