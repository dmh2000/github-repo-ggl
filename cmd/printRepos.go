package cmd

import "fmt"

func printRepos(names []string, stars []int) {
	for i, name := range names {
		if i < 9 {
			fmt.Printf("%d. %*s - %d\n", i+1, -33, name, stars[i])
		} else {
			fmt.Printf("%d. %*s - %d\n", i+1, -32, name, stars[i])
		}
	}
}