package docker

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Registry interface {
	Names() []string

	ApiUrl() string
	TokenUrl(repo string) string

	Transport(ctx context.Context, repo string) (http.RoundTripper, error)

	NormalizeRepo(repo string) string
	GetRepoUrl(repo string) string
	GetTagUrl(repo, tag string) string
	GetOwner(repo string) string
}

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
