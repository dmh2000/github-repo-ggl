/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-repo",
	Short: "Example app for fetching repositories from GitHub using the GitHub GraphQL API",
	Long: `This is an example app for fetching repositories from GitHub using the GitHub GraphQL API.
It has 3 subcommands:
raw      [owner] - print list of repositories of the owner with stars count using the raw client
shurcool [owner] - print list of repositories of the owner with stars count using the shurcool client
goog     [owner] - print list of repositories of the owner with stars count using the google/go-github client
`,
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}


