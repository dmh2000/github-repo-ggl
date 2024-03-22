# Find Out Everything About a GitHub Repo With The GitHub GraphQL API And Go

## Introduction

A lot of us are used to interacting with our github repos on the command line and in the GitHub web app. More advanced users can automate their interactions using the GitHub API's. GitHub has two API's: REST and GraphQL. Either one lets you access and automate processes. According to the documentation, the GraphQL API is targeted at more advanced usage.

This article is about using the GraphQL API to look at the details about any public repository, using Go.

### REST API

From [GitHub REST API Docs](https://docs.github.com/en/rest): "You can use GitHub's API to build scripts and applications that automate processes, integrate with GitHub, and extend GitHub. For example, you could use the API to triage issues, build an analytics dashboard, or manage releases."

### GraphQL API

From [GitHub GraphQL API Docs](https://docs.github.com/en/graphql): "To create integrations, retrieve data, and automate your workflows, use the GitHub GraphQL API. The GitHub GraphQL API offers more precise and flexible queries than the GitHub REST API."

### Differences

[REST vs GraphQL](https://aws.amazon.com/compare/the-difference-between-graphql-and-rest/)

"Under REST architecture, data is returned to the client from the server in the whole-of-resource structure specified by the server."

"A data format describes how you would like the server to return the data, including objects and fields that match the server-side schema"

On other words, GraphQL servers and client can mix data from multiple resources and specify just what is needed, where REST sends a single 'document'.

### Which One to Use

Get the scoop from the source:

[Comparing GitHub's REST API and GraphQL API](https://docs.github.com/en/rest/about-the-rest-api/comparing-githubs-rest-api-and-graphql-api)

In short, the GraphQL API allows fine grained access to its resources, where the REST API is less flexible and may give you more information than you might want. That's basically the difference between REST and GraphQL. HOWEVER, the easiest option to query the GraphQL API doesn't provide that level of granularity. More on that below.

## Three ways to query the GraphQL API

Unlike the REST API, all access to the GitHub GraphQL API require authentication. Examples of that will be in the code.

1. Raw POST Requests

At the low level, a client queries a GraphQL API using an HTTP POST request. The payload is a GraphQL formatted structure that specifies what you want to get. This is dooable, but can be kind of klunky and hard coded. It's possible to handcraft the POST request payload, but it is a bit tricky to get everything right. So most (all?) users will use a GraphQL Client package to simplify the process. There is an example below.

2. Google go-github

Google has created a Go package that supports accessing the GitHub GraphQL API, at [google/go-github](https://github.com/google/go-github). This is the easiest way to get at the GraphQL API, because it takes care of all the underlying GraphQL magic. It provides types and methods that correspond to the REST API.

This is the easiest way to go (pun), but the drawback is that is works like the REST API, returns whole documents rather than more fine grained requests that GraphQL is about.

3. A GraphQL client

For more fine grained access but with easier code than a raw POST, you can use a full client package. [Here's a list of libraries for Go](https://graphql.org/code/#go). Scroll down for clients. There are a couple of clients that have at least 1K GitHub stars. I chose [shurcooL/graphql](https://github.com/shurcooL/graphql), pretty easy to use, sort of.

With a direct client, there is more work setting up types to match the requests. But it allows full up GraphQL queries that can drill down to exactly what you want.

Since the GraphQL query specifies the elements to be returned for each repository, the result is an array of 3 members. The total bytes in my test was on the order of 1K bytes.

### Size of Returned Data

Targeting the 'octocat' repo, the size of the returned data in bytes was calculated.

1. raw
   Since the raw GraphQL query specifies the elements to be returned for each repository, the result is an array of 3 members. The total bytes in my test was on the order of 1K bytes.

2. go-github
   This one returns an array of Repository type. That type has over 100 members, even though the request only wants 3 of them. The total bytes in my test was on the order of 26K bytes.

3. shurcool
   Since the shurcool client GraphQL query specifies the elements to be returned for each repository, the result is an array of 3 members. The total bytes in my test was on the order of 1K bytes.

Conclusion: using the go-github API is convenient, if there is any concern over the overhead of the returned data, it may be worthwhile to use a more granular method. Probably using the shurcool or other GraphQL direct client.

## GitHub GraphQL API Schema

Figuring out what you want to do is a bit daunting. The GitHub GraphQL API schema is massive, a bit over 1MB in size. [You can download the schema](https://docs.github.com/en/graphql/overview/public-schema). But the schema is sort of opaque to someone (like me) who isn't a GraphQL or GitHub API expert. You can grep it to find what you are looking for. A simpler way can be to use the [GitHub GraphQL Explorer](https://docs.github.com/en/graphql/overview/public-schema). Pretty easy to search, even if you only have a vague idea of what you are looking for. And it then gives you a definition of the type layout you need to set up.

In this case I first searched the schema for 'repositories' and found quite a few hits. It wasn't clear to me what to use. I went to the Explorer and it was much more friendly. It took me a bit of flailing but I found the 'search' query which is what you want to find specific things. "Perform a search across resources, returning a maximum of 1,000 results". In this case I wanted a list of public repositories from any owner.

The query looks like this, with the variable 'queryString'. 'org' is the name of the account owner you want to search. In this case I only wanted the repo names from the type 'Repository'. That type has a bunch of data you can add to it. Go into the explorer, search 'Docs' and look for 'Repository', no prefix or suffix. That gives you the entire type definition. In most cases you will use that type definition to create or use a corresponding Go struct type.

With this query I was able to get a list for any owner account. Substitute whichever owner in place of 'octocat' in the queryString. Add any of the fields defined in the Repository type. If you have a GitHub account, login and go to the Explorer and try this out.

```graphql
query ($queryString: String!) {
  search(query: $queryString, type: REPOSITORY, first: 100) {
    repositoryCount
    edges {
      node {
        ... on Repository {
          name
        }
      }
    }
  }
}

{
"queryString": "owner:octocat"
}

```

## The Example Application

The code of this example app builds a command line utility that lets you query for information about any public github repo. The app has examples of the three types of access that fetches a list of repositories for a specified org/user.

The app is called 'gh-repo', with 3 subcommands: raw,go-github,shurcool and one argument 'user'.

- $gh-schema raw [user]
- $gh-schema go-github [user]
- $gh-schema shurcool [user]

It does the access using your account based on the GITHUB_TOKEN environment variable. The app has code that does the authentication.

The API used in the example code is the top level "query search" which can be used to search for several different GitHub objects. In this case a search of repositories for a specified owner is performed. You can find its definition in the GraphQL Explorer. Open the Schema Docs (file icon in upper left of the dialog), click on "Query", then scroll down to "search", or use hourglass to find "search".

There are other ways to list particular things that give a similar result. You can use the GraphQL Explorer or go into the google/go-github reference docs and look for matching types such as "RepositoriesService".

## The Code

Important! I use GitHub Copilot with VS Code. Without Copilot, it would have taken me 10 times longer to figure out exactly what to do. I use Copilot all the time but this case really made it worth the $10 a month. Not surprisingly, Copilot knows about the GraphQL types and really filled in a lot of the type information.

Each of the approaches: raw,go-github and shurcool perform the repo search and return the repo name,ID and stars count.

All access requires an authenticated client one way or another. The example code has the auth code and expects the environment variable GITHUB_TOKEN.

### Deriving Types

For the raw and shurcooL code, you will need to derive some types that match the GitHub schema. If you use the go-github package, it will have all the types defined. For the other two procedures, you can get the type specs from the GitHub schema, or easier from their GraphQL Explorer.

For the 'search' operation, do this:

- open the Explorer
- clock on the file icon in the upper left of the dialog
- search for Query and click on it
- scroll down

### RAW

See pkg/raw/raw.go.

This approach cobbles up a POST request in GraphQL language, sends it to the API and gets a JSON result back. To unmarshall the JSON, a type must be set up that matches the return format. The code defines a type named repoData that mirrors the JSON structure and has the appropriate annotations.

- define the query string
- configure the return type struct
- get an HTTP client
- concat the query strings
- create the POST request with authentication
- send the POSt
- get the JSON results
- unmarshall it
- extract and return the data

See the more detailed explanation generated by GitHub Copilot in the [appendix](#pkgrawrawgo)

### Go-GitHub

See pkg/goog/goog.go

This approach is easiest. Just look up the desired type and methods in the go-github docs and code it up.

- get an authenticated client : github.NewClient(nil).WithAuthToken
- perform the search : client.Search.Repositories
- extract the data and return it

See the more detailed explanation generated by GitHub Copilot in the [appendix](#pkggooggooggo)

### ShurcooL

See pkg/ghshurcool/ghshurcool.go.

This approach is in the middle between raw and go-github. It does require defining the query-search types. On the other hand, it allows fine grained access that the go-github API
doesn't allow.

- configure the result types
- get an authenticated OAuth client
- wrap it with the GraphQL client
- configure the query string
- make the request
- extract the data and return it

Note: The shurcooL/graphql package provides 'graphql' types including String, Float ID, Int and Int32. However, in the source comments and issues it says that the graphql types are not need. I had no problem using string, int, float directly.

## Setup And Run

### Auth Token

- You need a GitHub account.
- In order to access the GraphQL API, you need to get an auth token from your GitHub Account. The API requires authentication (The REST API does not). To get one, do the following:

  - login to your account
  - in the upper right corner, click on your profile picture
  - in the drop down list select "Settings"
  - on the Settings page, select "Developer Settings".
  - on the left menu, select "Personal Access Tokens"
  - select "Fine-grained tokens". (These allow more targeted permissions on the token)
  - generate a new token
    - by default, new tokens allow only read-only access to the API. That's enough for this exercise.
    - it may ask for your password or passkey to verify it's you
    - follow the instructions. Remember that once you generate the token, you need to copy it and store it somewhere (NOT IN A REPO). You won't be able to look at it again on the GitHub page.
    - For this exercise, the code will look for an environment variable named "GITHUB_TOKEN".
    - The token will look something like this "github_pat_blahblahblah..."

### Install Go

The code here requires Go 1.18 or later (for go.work and module mode). If you are reading this, you probably already have Go, but if not, go to [Go Download and Install](https://go.dev/doc/install).

Note : All the code here is built and tested on Linux Mint 21.2 Victoria. Compatible with Ubuntu 22.04LTA. But Go on any platform should work.

### Clone The Repo

If you haven't already.

[The repo with code](https://github.com/dmh2000/github-repo-ggl)

### Run It

- printenv | grep GITHUB_TOKEN (just checking)

- $cd into the repo (top level)
- $go run . raw [owner name]
- $go run . goog [owner name]
- $go run . shurcool [owner name]

OR

- $cd into the repo (top level)
- $source test.sh (runs all three against the 'octocat' owner)

## Appendix

The following are descriptions, generated by GitHub Copilot, of the three procedures to fetch the repos. You don't need to read this part unless you are interested in what Copilot can do for documentation.

### pkg/raw/raw.go

The provided Go code is used to fetch repository data from GitHub's GraphQL API.

The query variable is an array of strings that, when concatenated, form a GraphQL query. This query is used to fetch repositories owned by a specific user. The owner's name is inserted into the query at the third index of the array.

The repoData struct is used to unmarshal the JSON response from the GitHub API. It has nested structs that mirror the structure of the JSON response. Each Node struct represents a repository and contains fields for the repository's name, ID, and star count.

The FetchRepos function is the main function in this code. It takes an owner's name as an argument and returns lists of repository IDs, names, and star counts, or an error if something goes wrong.

The function starts by getting the GitHub API token from the environment variables. It then creates an HTTP client with a 10-second timeout.

Next, it inserts the owner's name into the GraphQL query and joins the query strings into a single string. It creates a new HTTP request with this query as the body. The request is a POST request to the GitHub GraphQL API endpoint. The function adds necessary headers to the request, including the Authorization header with the API token and the Content-Type header set to "application/json".

The function then sends the request using the HTTP client and checks for any errors. If the request is successful, it reads the entire response body and unmarshals the JSON into the repoData struct.

Finally, the function iterates over the Edges slice in the repoData struct, which contains the data for each repository. It appends the ID, name, and star count of each repository to separate slices. These slices are then returned by the function.

### pkg/goog/goog.go

The provided Go code is a function named FetchRepos that fetches repositories from GitHub for a given user or organization. The function takes a string argument owner which is the username of the GitHub user or organization. It returns four values: slices of repository IDs, names, and star counts, and an error if something goes wrong.

The function starts by creating a new GitHub client using the github.NewClient(nil) function. It then authenticates the client with a GitHub token fetched from the environment variables using os.Getenv("GITHUB_TOKEN).

Next, it creates a github.SearchOptions struct and a query string that specifies the user or organization to fetch repositories from. The query string is formatted using fmt.Sprintf("user:%s", owner), which replaces %s with the owner string.

The function then calls client.Search.Repositories(context.Background(), query, opt) to send a request to GitHub's API to search for repositories. This function takes three arguments: a context (in this case, a background context), the query string, and the search options. It returns a RepositoriesSearchResult struct, a Response struct, and an error.

If an error occurs during the search, the function prints the response and error, and returns nil for the repository IDs, names, and star counts, along with the error.

If the search is successful, the function assigns the repositories from the search result to the repos variable. It then initializes three slices to hold the repository IDs, names, and star counts.

The function then iterates over the repos slice. For each repository, it appends the repository's ID, name, and star count to the respective slices. These slices are then returned by the function, along with nil for the error.

### pkg/ghshurcool/ghshurcool

This code is designed to interact with the GitHub GraphQL API to fetch repository information.

The Repository struct is a representation of a GitHub repository. It contains fields for the repository's name, ID, and the total count of its stargazers. The graphql tags indicate how these fields map to the corresponding fields in the GraphQL response.

The RepositoryEdge struct represents an edge in a GraphQL connection. It includes a cursor for pagination and a Node of type Repository.

The PageInfo struct represents information about pagination in a GraphQL connection. It includes fields for the start and end cursors, and boolean fields indicating whether there are more pages to fetch in either direction.

The RepositoryConnection struct represents a connection in GraphQL. It includes a slice of edges, each containing a RepositoryEdge, a slice of nodes, each containing a Repository, and a PageInfo object.

The repos variable is a struct that represents the structure of the expected GraphQL response. It includes a Search field, which is a struct that represents a search operation in the GitHub GraphQL API.

The variables map is used to store variables that will be used in the GraphQL query.

The FetchRepos function is used to fetch repositories from the GitHub GraphQL API. It takes an owner string as an argument, which is used to construct the search query. It uses the oauth2 package to authenticate the request with a GitHub token, which is fetched from the environment variables. It then creates a new GraphQL client and executes the query. If the query is successful, it iterates over the edges in the response, extracting the name, ID, and stargazer count of each repository, and returns these as slices. If the query fails, it returns an error.
