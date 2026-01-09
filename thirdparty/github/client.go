package github

import (
	"Blog-Backend/consts"
	"context"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	cli *githubv4.Client
}

func NewClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(consts.EnvDiscussionToken)},
	)
	// 初始化httpclient
	httpClient := oauth2.NewClient(context.Background(), src)
	cli := githubv4.NewClient(httpClient)
	return cli
}
