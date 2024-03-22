package goog

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v60/github"
)


func FetchRepos(owner string) ([]string, []string, []int, error) {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	opt := &github.SearchOptions{}
	query := fmt.Sprintf("user:%s", owner)
	result,response,err := client.Search.Repositories(context.Background(), query, opt)
	if err != nil {
		fmt.Println(response)
		fmt.Println(err)
		return  nil,nil,nil,err
	}

	repos := result.Repositories

	// uncomment these lines to see the size of the returned data
	// buf := bytes.Buffer{}
	// enc := gob.NewEncoder(&buf)
	// err = enc.Encode(repos)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
	// fmt.Printf("Size of go-github: %d\n", buf.Len())

	var names []string
	var stars []int
	var ids []string
	for _, repo := range repos {
		ids = append(ids, *repo.NodeID)
		names = append(names, *repo.Name)
		stars = append(stars, *repo.StargazersCount)
	}
	return ids, names, stars, nil
}
