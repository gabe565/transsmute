package docker

import (
	"net/http"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.String("ghcr-username", "", "GitHub username for ghcr.io")
	if err := viper.BindPFlag("ghcr.username", flag.Lookup("ghcr-username")); err != nil {
		panic(err)
	}

	flag.String("ghcr-password", "", "GitHub personal access token for ghcr.io")
	if err := viper.BindPFlag("ghcr.password", flag.Lookup("ghcr-password")); err != nil {
		panic(err)
	}
}

type Ghcr struct{}

func (g Ghcr) Name() string {
	return "ghcr.io"
}

func (g Ghcr) ApiUrl() string {
	return "https://ghcr.io"
}

func (g Ghcr) TokenUrl(repo string) string {
	return g.ApiUrl() + "/token?service=ghcr.io&scope=repository:" + repo + ":pull"
}

func (g Ghcr) Transport(repo string) http.RoundTripper {
	username := viper.GetString("ghcr.username")
	password := viper.GetString("ghcr.password")

	return registry.WrapTransport(
		http.DefaultTransport,
		g.TokenUrl(repo),
		username,
		password,
	)
}

func (g Ghcr) NormalizeRepo(repo string) string {
	return repo
}

func (g Ghcr) GetRepoUrl(repo string) string {
	return "https://" + repo
}

func (g Ghcr) GetTagUrl(repo, tag string) string {
	return g.GetRepoUrl(repo)
}

func (g Ghcr) GetOwner(repo string) string {
	return strings.SplitN(repo, "/", 3)[1]
}
