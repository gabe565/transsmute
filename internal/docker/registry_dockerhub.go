package docker

import (
	"context"
	"strings"

	"github.com/gabe565/transsmute/internal/config"
	"github.com/google/go-containerregistry/pkg/authn"
)

type DockerHub struct { //nolint:revive
	username, password string
}

func NewDockerHub(conf config.DockerHub) (*DockerHub, error) {
	return &DockerHub{
		username: conf.Username,
		password: conf.Password,
	}, nil
}

func (d DockerHub) Names() []string {
	return []string{"", "docker.io"}
}

func (d DockerHub) APIURL() string {
	return "https://registry.hub.docker.com"
}

func (d DockerHub) TokenURL(repo string) string {
	return "https://auth.docker.io/token?service=registry.hub.docker.com&scope=repository:" + repo + ":pull"
}

func (d DockerHub) Authenticator(_ context.Context, _ string) (authn.Authenticator, error) {
	return &authn.Basic{
		Username: d.username,
		Password: d.password,
	}, nil
}

func (d DockerHub) NormalizeRepo(repo string) string {
	repo = strings.TrimPrefix(repo, "docker.io/")
	if !strings.Contains(repo, "/") {
		return "library/" + repo
	}
	return repo
}

func (d DockerHub) GetRepoURL(repo string) string {
	if strings.HasPrefix(repo, "library/") {
		return "https://hub.docker.com/_/" + strings.TrimPrefix(repo, "library/")
	}
	return "https://hub.docker.com/r/" + repo
}

func (d DockerHub) GetTagURL(repo, tag string) string {
	return d.GetRepoURL(repo) + "/tags?name=" + tag
}

func (d DockerHub) GetOwner(repo string) string {
	return strings.SplitN(repo, "/", 2)[0]
}
