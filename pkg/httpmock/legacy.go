package httpmock

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// TODO: clean up methods in this file when there are no more callers

func (r *Registry) StubResponse(status int, body io.Reader) {
	r.Register(MatchAny, func(req *http.Request) (*http.Response, error) {
		return httpResponse(status, req, body), nil
	})
}

func (r *Registry) StubWithFixturePath(status int, fixturePath string) func() {
	fixtureFile, err := os.Open(fixturePath)
	r.Register(MatchAny, func(req *http.Request) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		return httpResponse(200, req, fixtureFile), nil
	})
	return func() {
		if err == nil {
			fixtureFile.Close()
		}
	}
}

func (r *Registry) StubRepoInfoResponse(owner, repo, branch string) {
	r.Register(
		GraphQL(`query RepositoryInfo\b`),
		StringResponse(fmt.Sprintf(`
		{ "data": { "repository": {
			"id": "REPOID",
			"name": "%s",
			"owner": {"login": "%s"},
			"description": "",
			"defaultBranchRef": {"name": "%s"},
			"hasIssuesEnabled": true,
			"viewerPermission": "WRITE"
		} } }
		`, repo, owner, branch)))
}

func (r *Registry) StubRepoResponse(owner, repo string) {
	r.StubRepoResponseWithPermission(owner, repo, "WRITE")
}

func (r *Registry) StubRepoResponseWithPermission(owner, repo, permission string) {
	r.Register(GraphQL(`query RepositoryNetwork\b`), StringResponse(RepoNetworkStubResponse(owner, repo, "master", permission)))
}

func RepoNetworkStubResponse(owner, repo, defaultBranch, permission string) string {
	return fmt.Sprintf(`
		{ "data": { "repo_000": {
			"id": "REPOID",
			"name": "%s",
			"owner": {"login": "%s"},
			"defaultBranchRef": {
				"name": "%s"
			},
			"viewerPermission": "%s"
		} } }
	`, repo, owner, defaultBranch, permission)
}
