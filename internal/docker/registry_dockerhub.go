package docker

import (
	"github.com/heroku/docker-registry-client/registry"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func init() {
	flag.String("dockerhub-username", "", "DockerHub username for private repos")
	if err := viper.BindPFlag("dockerhub.username", flag.Lookup("dockerhub-username")); err != nil {
		panic(err)
	}

	flag.String("dockerhub-password", "", "DockerHub password for private repos")
	if err := viper.BindPFlag("dockerhub.password", flag.Lookup("dockerhub-password")); err != nil {
		panic(err)
	}
}

type DockerHub struct{}

func (d DockerHub) Name() string {
	return ""
}

func (d DockerHub) ApiUrl() string {
	return "https://registry.hub.docker.com"
}

func (d DockerHub) TokenUrl(repo string) string {
	return "https://auth.docker.io/token?service=registry.hub.docker.com&scope=repository:" + repo + ":pull"
}

func (d DockerHub) Transport(repo string) http.RoundTripper {
	return registry.WrapTransport(
		http.DefaultTransport,
		d.TokenUrl(repo),
		viper.GetString("dockerhub.username"),
		viper.GetString("dockerhub-password"),
	)
}

func (d DockerHub) NormalizeRepo(repo string) string {
	if !strings.Contains(repo, "/") {
		return "library/" + repo
	}
	return repo
}

func (d DockerHub) GetRepoUrl(repo string) string {
	if strings.HasPrefix(repo, "library/") {
		return "https://hub.docker.com/_/" + strings.TrimPrefix(repo, "library/")
	}
	return "https://hub.docker.com/r/" + repo
}

func (d DockerHub) GetTagUrl(repo, tag string) string {
	return d.GetRepoUrl(repo) + "/tags?name=" + tag
}

func (d DockerHub) GetOwner(repo string) string {
	return strings.SplitN(repo, "/", 2)[0]
}
