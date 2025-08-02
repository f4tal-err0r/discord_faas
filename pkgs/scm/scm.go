package scm

import (
	"fmt"
	"net/url"
)

type Client struct {
	Repos map[string]Repofuncs
	SCMClient
}

type SCMClient interface {
	GetFuncs() ([]Repofuncs, error)
	GetSha() string
	GetBranch() string
	FetchCurrenRepo() error
	IsCurrent() (bool, error)
}

type Repofuncs struct {
	Name    string
	Path    string
	RepoURL *url.URL
	Branch  string
}

// TODO: Support private repos via oauth
func NewClient(repoURL string) (*Client, error) {
	url, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}

	var scm SCMClient

	switch url.Host {
	case "github.com":
		//TODO
	default:
		return nil, fmt.Errorf("unsupported host: %s", url.Host)
	}
	c := &Client{
		Repos: make(map[string]Repofuncs),
	}
	return c, nil
}

func (c *Client) GetFileTree() string {
	return ""
}
