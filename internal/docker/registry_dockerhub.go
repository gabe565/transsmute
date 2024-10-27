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

func (d DockerHub) Names() []string {
	return []string{"", "docker.io"}
}

func (d DockerHub) APIURL() *url.URL {
	return &url.URL{Scheme: "https", Host: "registry.hub.docker.com"}
}

func (d DockerHub) TokenURL(repo string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "auth.docker.io",
		Path:   "/token",
		RawQuery: url.Values{
			"service": []string{"registry.hub.docker.com"},
			"scope":   []string{"repository:" + repo + ":pull"},
		}.Encode(),
	}
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

func (d DockerHub) GetRepoURL(repo string) *url.URL {
	u := &url.URL{Scheme: "https", Host: "hub.docker.com"}
	if strings.HasPrefix(repo, "library/") {
		u.Path = path.Join(u.Path, "_", strings.TrimPrefix(repo, "library/"))
	} else {
		u.Path = path.Join(u.Path, "r", repo)
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
	return strings.SplitN(repo, "/", 2)[0]
}
