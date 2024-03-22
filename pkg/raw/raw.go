package raw

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// formatting the graphql json is tricky
var query = []string{
	"{\"query\":",
	`"query {search(query: \"owner:`,
	"",
	// owner name
	`\", type: REPOSITORY, first: 100) { repositoryCount edges { node { ... on Repository { name id stargazerCount}}}}}"`,
	"}",
}

// JSON result format (figured out manually by printing the result types)
// probably should use a client package
type repoData struct {
	Data struct {
		Search struct {
			Edges []struct {
				Node struct {
					Name           string `json:"name"`
					ID             string `json:"id"`
					StargazerCount int    `json:"stargazerCount"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"search"`
	} `json:"data"`
}

func FetchRepos(owner string) ([]string, []string, []int, error) {
	token := os.Getenv("GITHUB_TOKEN")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// construct query
	query[2] = owner
	q := strings.Join(query,"")

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", nil)
	if err != nil {
		return nil, nil, nil, err
	}

	reqBody := io.NopCloser(strings.NewReader(q))

	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Body = reqBody
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil,nil, nil, err
	}

	var repos repoData
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, nil, nil, err
	}

	// uncomment these lines to see the size of the returned data
	// buf := bytes.Buffer{}
	// enc := gob.NewEncoder(&buf)
	// err = enc.Encode(repos)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
	// fmt.Printf("Size of raw: %d\n", buf.Len())

	ids := []string{}
	names := []string{}
	stars := []int{}

	for _, v := range repos.Data.Search.Edges {
		ids = append(ids, v.Node.ID)
		names = append(names, v.Node.Name)
		stars = append(stars, v.Node.StargazerCount)
	}

	return ids, names, stars, nil
}
