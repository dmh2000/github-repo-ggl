#!/bin/sh

echo "# raw POST         =================================================="
go run . raw octocat
echo "# google/go-github =================================================="
go run . goog octocat
echo "# shurcool/graphql ================================================="
go run . shurcool octocat