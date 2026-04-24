package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	region  string
	profile string
	interval int
)

var rootCmd = &cobra.Command{
	Use:   "stackwatch",
	Short: "Monitor CloudFormation and Terraform stack drift in real time",
	Long: `stackwatch is a lightweight CLI tool that monitors your infrastructure
stacks for configuration drift, alerting you when resources deviate
from their expected state.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS region to target")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS profile to use")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 60, "Polling interval in seconds")
}
