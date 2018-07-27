package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/hryyan/b2"
	"github.com/spf13/cobra"
)

var publicBucketCmd = &cobra.Command{
	Use:   "public [bucket ..]",
	Short: "Public bucket",
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
					bucket.BucketType = b2.PUBLIC
					_, err := client.UpdateBucket(bucket, false)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(B2_LIBRARY_ERROR_EXIT)
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

var privateBucket = &cobra.Command{
	Use:   "private [bucket ..]",
	Short: "Private bucket",
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
					bucket.BucketType = b2.PRIVATE
					_, err := client.UpdateBucket(bucket, false)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(OPERATION_ERROR_EXIT)
					} else {
						fmt.Printf("Private bucket %s successed!\n", name)
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

var flushAllUnfinishPartCmd = &cobra.Command{
	Use:   "flush [bucket..]",
	Short: "Flush all unfinished part",
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
					files, err := client.ListUnfinishedLargeFiles(
						bucket.BucketId,
						"",
						"",
						100,
					)
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(B2_LIBRARY_ERROR_EXIT)
					}

					for _, file := range files {
						err = client.CancelLargeFile(file.FileId)
						if err != nil {
							fmt.Println(err.Error())
							os.Exit(B2_LIBRARY_ERROR_EXIT)
						}
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

var public bool = false

var newBucketCmd = &cobra.Command{
	Use:   "bucket [bucket ..]",
	Short: "New bucket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := login()
		bucketType := ""

		if public {
			bucketType = b2.PUBLIC
		} else {
			bucketType = b2.PRIVATE
		}

		for _, name := range args {
			_, err := client.CreateBucket(name, bucketType,
				map[string]string{}, []b2.CorsRule{}, []b2.LifecycleRule{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(OPERATION_ERROR_EXIT)
			}
		}
	},
}

var newCmd = &cobra.Command{
	Use:   "new command",
	Short: "New bucket",
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register a b2 account",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("open", "https://www.backblaze.com/b2/sign-up.html")
			cmd.Run()
		} else {
			fmt.Println("https://www.backblaze.com/b2/sign-up.html")
		}
	},
}

func init() {
	newBucketCmd.Flags().BoolVarP(
		&public,
		"public",
		"p",
		false,
		"public bucket")

	newCmd.AddCommand(newBucketCmd)

	rootCmd.AddCommand(publicBucketCmd, privateBucket, flushAllUnfinishPartCmd, newCmd, registerCmd)
}
