package github

import "github.com/shurcooL/githubv4"

type Client struct {
	cli *githubv4.Client
}

func NewClient(cli *githubv4.Client) *Client {
	return &Client{cli: cli}
}
