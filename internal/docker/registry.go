package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
)

type Registry interface {
	Names() []string

	APIURL() string
	TokenURL(repo string) string

	Authenticator(ctx context.Context, repo string) (authn.Authenticator, error)

	NormalizeRepo(repo string) string
	GetRepoURL(repo string) string
	GetTagURL(repo, tag string) string
	GetOwner(repo string) string
}

//nolint:gochecknoglobals
var registries []Registry

func SetupRegistries() error {
	ghcr, err := NewGhcr()
	if err != nil {
		return err
	}

	registries = []Registry{
		ghcr,
		&DockerHub{},
	}
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
