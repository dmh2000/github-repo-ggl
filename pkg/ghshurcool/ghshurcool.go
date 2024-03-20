package ghshurcool

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Repository struct {
	Name string `graphql:"name"`
	ID   string `graphql:"id"`
	StarGazers struct {
		TotalCount int `graphql:"totalCount"`
	} `graphql:"stargazers"`
} // @name Node

type RepositoryEdge struct {
	Cursor string `graphql:"cursor"`
	Node Repository `graphql:"node"`
} // @name RepositoryEdge

type PageInfo struct {
	EndCursor   string `graphql:"endCursor"`
	HasNextPage bool   `graphql:"hasNextPage"`
	HasPreviousPage bool `graphql:"hasPreviousPage"`
	StartCursor string `graphql:"startCursor"`
} // @name PageInfo

type RepositoryConnection struct {
	RepositoryConnection struct {
		Edges []struct {
			RepositoryEdge RepositoryEdge `graphql:"node"`
		} `graphql:"edges"`
		Nodes []struct {
			Node Repository `graphql:"node"`
		} `graphql:"nodes"`
		PageInfo PageInfo `graphql:"pageInfo"`
	} `graphql:"repositoryConnection(first: 10, after: $cursor)"`
} // @name RepositoryConnection

var repos struct {
	Search struct {
		Edges []struct {
			Node struct {
				SearchedRepository struct {
					Name        string `graphql:"name"`
					ID 	   		string `graphql:"id"`
					StarGazers struct {
						TotalCount int `graphql:"totalCount"`
					} `graphql:"stargazers"`
				} `graphql:"... on Repository"`
			}
		}
	} `graphql:"search(query: $query, type:REPOSITORY, first: $first)"`
}

var variables = map[string]any{}

//lint:ignore U1000 called by main
func FetchRepos(owner string) ([]string, []string, []int, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.github.com/graphql", httpClient)

	// update the variables
	variables["query"] = graphql.String(fmt.Sprintf("owner:%s", owner))
	variables["first"] = graphql.Int(100)

	// execute the query
	err := client.Query(context.Background(), &repos, variables)
	if err != nil {
		return nil, nil, nil, err
	}

	// uncomment this line to see the size of the returned data
	// fmt.Println(repos)

	var names []string
	var stars []int
	var ids []string
	for _, edge := range repos.Search.Edges {
		names = append(names, edge.Node.SearchedRepository.Name)
		stars = append(stars, edge.Node.SearchedRepository.StarGazers.TotalCount)
		ids = append(ids, edge.Node.SearchedRepository.ID)
	}
	return ids, names, stars, nil
}
