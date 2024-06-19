package docker

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gabe565/transsmute/internal/config"
	"github.com/google/go-containerregistry/pkg/authn"
)

type Registry interface {
	Names() []string

	APIURL() *url.URL
	TokenURL(repo string) *url.URL

	Authenticator(ctx context.Context, repo string) (authn.Authenticator, error)

	NormalizeRepo(repo string) string
	GetRepoURL(repo string) *url.URL
	GetTagURL(repo, tag string) *url.URL
	GetOwner(repo string) string
}

//nolint:gochecknoglobals
var registries []Registry

func SetupRegistries(conf config.Docker) error {
	ghcr, err := NewGhcr(conf.GHCR)
	if err != nil {
		return err
	}

	dockerhub, err := NewDockerHub(conf.DockerHub)
	if err != nil {
		return err
	}

	registries = []Registry{ghcr, dockerhub}
	return nil
}

var ErrInvalidRegistry = errors.New("no registry for repo")

func FindRegistry(repo string) (Registry, error) {
	for _, registry := range registries {
		for _, name := range registry.Names() {
			if strings.HasPrefix(repo, name) {
				return registry, nil
			}
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidRegistry, repo)
}
