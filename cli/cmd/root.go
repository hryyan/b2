package cmd

import (
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/hryyan/b2"
)

var (
	verbose     bool   = false
	sessionPath string = ""
)

var rootCmd = &cobra.Command{
	Use: "b2",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}

func init() {
	cobra.OnInitialize(initSession)
	viper.BindEnv("B2_ACCOUNT_ID")
	viper.BindEnv("B2_APPLICATION_KEY")

	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"Produce verbose output")

}

func initSession() {
	home, err := homedir.Dir()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	sessionPath = path.Join(home, ".b2_session")

	if _, err = os.Stat(sessionPath); os.IsNotExist(err) {
		if _, err := os.Create(sessionPath); err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("Save session in %s\n", sessionPath)
		}
	}

	_, err = os.Open(sessionPath)
	if err != nil {
		log.Println(err.Error())
	}
}
