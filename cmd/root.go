package cmd

import (
	"fmt"
	"github.com/redhat-appstudio/qe-tools/cmd/estimate"
	"github.com/redhat-appstudio/qe-tools/cmd/webhook"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redhat-appstudio/qe-tools/cmd/coffeebreak"
	"github.com/redhat-appstudio/qe-tools/cmd/prowjob"
	"github.com/redhat-appstudio/qe-tools/cmd/sendslackmessage"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qe-tools",
	Short: "The CLI containing useful tools used by RHTAP QE",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qe-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(prowjob.ProwjobCmd)
	rootCmd.AddCommand(coffeebreak.CoffeeBreakCmd)
	rootCmd.AddCommand(sendslackmessage.SendSlackMessageCmd)
	rootCmd.AddCommand(webhook.WebhookCmd)
	rootCmd.AddCommand(estimate.EstimateTimeToReviewCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
