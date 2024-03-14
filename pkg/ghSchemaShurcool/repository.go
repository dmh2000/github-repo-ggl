package Shurcool

import (
	"fmt"
)

/*
mutation createRepository($input: CreateRepositoryInput!) {
  createRepository(input: $input) {
    clientMutationId
  }
}
Variables
{
  "input" : {
    "name": "repository_name"
  	"ownerId": "",
	"description": "This is your first repository",
    "visibility": "PUBLIC","PRIVATE","INTERNAL"
	"template": false,
	"homepageUrl": "https://github.com",
	"hasWikiEnabled": false,
	"hasIssuesEnabled": true,
	"teamId": ""
	}
}

type CreateRepository struct {
	CreateRepository struct {
			ClientMutationID string `graphql:"clientMutationId"`
	} `graphql:"createRepository(input:{name:\"newrepo\" visibility:PUBLIC ownerId: \"ownerid\" description:\"This is your first repository\" template:false homepageUrl:\"https://github.com\" hasWikiEnabled:false hasIssuesEnabled:true teamId:\"teamid})"`
} // @name CreateRepository
*/

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

func Run() {
	fmt.Println("repository")
}