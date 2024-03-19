# Find Out Everything About A GitHub Repo With The GitHub GraphQl API (With Golang)

## Introduction

A lot of us are used to interacting with our github repos on the command line and in the GitHub web app. More advanced users can automate their interactions using the GitHub API's. GitHub has two API's : REST and GraphQl. Either one lets you access and automate processes. According to the documentation, the GraphQl API is targeted at more advanced usage.

This article is about using the GraphQl API to look at the details about any public repository, using Go.

### REST API

From [Github REST API Docs](https://docs.github.com/en/rest) : "You can use GitHub's API to build scripts and applications that automate processes, integrate with GitHub, and extend GitHub. For example, you could use the API to triage issues, build an analytics dashboard, or manage releases."

### GraphQl API

From [GitHub GraphQl API Docs](https://docs.github.com/en/graphql) : "To create integrations, retrieve data, and automate your workflows, use the GitHub GraphQL API. The GitHub GraphQL API offers more precise and flexible queries than the GitHub REST API."

### Which One To Use

Get the scoop from the source:

[Comparing GitHub's REST API and GraphQL API](https://docs.github.com/en/rest/about-the-rest-api/comparing-githubs-rest-api-and-graphql-api)

In short, the GraphQl API allows fine grained access to its resources, where the REST API is less flexible and may give you more information than you might want. That's basically the difference between REST and GraphQl.

## The Application

The code of this app builds a command line utility that lets you query for information about any public github repo. It does the access using your account based on the GITHUB_TOKEN environment variable.

At the low level, querying a GraphQl API server requires sending a POST request to the server with a payload containing the proper GraphQl language. The response from a valid request will return a payload that has the requested data in JSON format. Or an error message.

Its possible to handcraft the POST request payload, but it is a bit tricky to get everything right. So most (all?) users will use a GraphQl Client package to simplify the process.

You can go to [Go Clients](https://graphql.org/code/#go) at the main GraphQL site. Scroll down to the client section. There are a couple of client packages that have at least 1K stars. I selected shurcool/graphql for this project. It worked pretty seemlessly.

## GitHub GraphQL API

The GitHub GraphQL API schema is massive, a bit over 1MB in size. [You can download the schema](https://docs.github.com/en/graphql/overview/public-schema). But, the schema is sort of opaque to someone (like me) who isn't a GraphQL or Github API expert. You can grep it to find what you are looking for. A simpler way can be to use the [BitHb GraphQL Explorer](https://docs.github.com/en/graphql/overview/public-schema). Pretty easy to search, even if you only have a vague idea of what you are looking for. And it then gives you a definition of the type layout you need to set up.

In this case I first searched the schema for 'repositories' and found quite a few hits. It wasn't clear to me what to use. I went to the Explorer and it was much more friendly. It took me a bit of flailing but I found the 'search' query which is what you want to find specific things. "Perform a search across resources, returning a maximum of 1,000 results". In this case I wanted a list of public repositories from any owner.

The query looks like this, with the variable 'queryString'. 'org' is the name of the account owner you want to search. In this case I only wanted the repo names from the type 'Repository'. That type has a bunch of data you can add to it. Go into the explorer, search 'Docs' and look for 'Repository', no prefix or suffix. That gives you the entire type definition.

With this query I was able to get a list for any owner account. Substitute whichever owner in place of 'octocat' in the queryString. Add any of the fields defined in the Repository type. If you have a GitHub account, go to the Explorer and try this out.

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
"queryString": "is:public archived:false org:octocat"
}

```

## The Code

Important! I use GitHub Copilot with VS Code. Without Copilot, it would have taken me 10 times longer to figure out exactly what to do. I use Copilot all the time but this case really made it worth the $10 a month. Not surprisingly, Copilot knows about the GraphQL types and really filled in a lot of the type information.

If you are familiar with marshall/unmarshall JSON types, the GraphQL types work similarly. You define a type, usually a struct, and annotate it's members with a 'graphql' tag. That tells Go how to map from the JSON data to Go types. ike this:

```go

var q struct {
	Human struct {
		Name   graphql.String
		Height graphql.Float  `graphql:"height(unit: METER)"`
	} `graphql:"human(id: \"1000\")"`
}

```

Note: The shurcooL package provides 'graphql' types including String, Float ID, Int and Int32. However in the source comments and issues it says that the graphql types are not need. I had no problem using string, int, float directly.

### Define the query type

### Define the result type(s)

### Execute the query

Now that we know what the query needs to look like, it needs to be translated to Go code. Using the shurcooL/graphql client, I was able to piece together the code to execute the query. The process looks something like this:

- Import "github.com/shurcooL/graphql" and "golang.org/x/oauth2"
- Get an authenticated http client with GITHUB_TOKEN, using OAuth.
- Create a graphql client and give it the http client
- Set up variables
- Execute the query
- Extract and return the requested data.

### Executing the Query

## Setup And Run

### Auth Token

- You need a GitHub account.
- In order to access the GraphQl API, you need to get an auth token from your GitHub Account. The API requires authentication (The REST API does not). To get one, do the following:

  - login to your account
  - in the upper right corner, click on your profile picture
  - in the drop down list select "Settings"
  - on the Settings page, select "Developer Settings".
  - on the left menu, select "Personal Access Tokens"
  - select "Fine-grained tokens". (These allow more targeted permissions on the token)
  - generate a new token
    - by default, new tokens allow only read-only access to the API. That's enough for this exercise.
    - it may ask for your password or passkey to verify its you
    - follow the instructions. Remember that once you generate the token, you need to copy it an store it somewhere (NOT IN A REPO). You won't be able to look at it again on the GitHub page.
    - For this exercise, the code will look for an environment variable named "GITHUB_TOKEN".
    - The token will look something like this "github_pat_blahblahblah..."

### Install Go

The code here requires Go 1.18 or later (for go.work and module mode). If you are reading this, you probably already have Go, but if not, go to [Go Download And Install](https://go.dev/doc/install).

Note : All the code here is built and tested on Linux Mint 21.2 Victoria. Compatible with Ubuntu 22.04LTA. But Go on any platform should work.

### Clone The Repo

If you haven't already.

[The repo with code](https://github.com/dmh2000/github-repo-ggl)

### Run It

- printenv | grep GITHUB_TOKEN (just checking)
- $cd cmd
- $go run .
