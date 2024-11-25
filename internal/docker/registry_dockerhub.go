package docker

import (
	"context"
	"net/url"
	"path"
	"strings"

	"gabe565.com/transsmute/internal/config"
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

func (d DockerHub) Match(repo string) bool {
	return strings.HasPrefix(repo, "docker.io/") || strings.Count(repo, "/") <= 1
}

func (d DockerHub) Authenticator(_ context.Context, _ string) (authn.Authenticator, error) {
	return &authn.Basic{
		Username: d.username,
		Password: d.password,
	}, nil
}

func (d DockerHub) GetRepoURL(repo string) *url.URL {
	u := &url.URL{Scheme: "https", Host: "hub.docker.com"}
	if strings.ContainsRune(repo, '/') {
		u.Path = path.Join(u.Path, "r", repo)
	} else {
		u.Path = path.Join(u.Path, "_", repo)
	}
	return u
}

func (d DockerHub) GetTagURL(repo, tag string) *url.URL {
	u := d.GetRepoURL(repo)
	u.Path = path.Join(u.Path, "tags")
	query := u.Query()
	query.Set("name", tag)
	u.RawQuery = query.Encode()
	return u
}

func (d DockerHub) GetOwner(repo string) string {
	if !strings.ContainsRune(repo, '/') {
		return "library"
	}
	return strings.SplitN(repo, "/", 2)[0]
}
