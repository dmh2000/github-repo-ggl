/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sqirvy.xyz/gh-repo/pkg/raw"
)

// shurcoolCmd represents the shurcool command
var rawCmd = &cobra.Command{
	Use:   "raw [owner]",
	Short: "print list of repositories of the owner with stars count",
	Long:  "Print list of repositories of the owner with stars count using a raw POST",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the owner name")
			return
		}
		owner := args[0]
		
		// see pkg/raw/raw.go
		_,names,stars,err := raw.FetchRepos(owner)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Check if the GITHUB_TOKEN is set")
			return
		}

		printRepos(names, stars)
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)
}
