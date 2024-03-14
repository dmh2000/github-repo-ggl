# Find Out Everything About A GitHub Repo With The GitHub GraphQl API (With Golang)

## Introduction

A lot of us are used to interacting with our github repos on the command line and in the GitHub web app. More advanced users can automate their interactions using the GitHub API's. GitHub has two API's : REST and GraphQl. Either one lets you access and automate processes. According to the documentation, the GraphQl API is targeted at more advanced usage.

This article is about using the GraphQl API to look at any public repository, using Go.

### REST API

From [Github REST API Docs](https://docs.github.com/en/rest) : "You can use GitHub's API to build scripts and applications that automate processes, integrate with GitHub, and extend GitHub. For example, you could use the API to triage issues, build an analytics dashboard, or manage releases."

### GraphQl API

From [GitHub GraphQl API Docs](https://docs.github.com/en/graphql) : "To create integrations, retrieve data, and automate your workflows, use the GitHub GraphQL API. The GitHub GraphQL API offers more precise and flexible queries than the GitHub REST API."

### Which One To Use

Get the scoop from the source:

[Comparing GitHub's REST API and GraphQL API](https://docs.github.com/en/rest/about-the-rest-api/comparing-githubs-rest-api-and-graphql-api)

In short, the GraphQl API allows fine grained access to its resources, where the REST API is less flexible and may give you more information than you might want. That's basically the difference between REST and GraphQl.

## Getting Ready

1. You need a GitHub account.
2. In order to access the GraphQl API, you need to get an auth token from your GitHub Account. The API requires authentication.
