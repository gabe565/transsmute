package docker

import (
	"context"
	"net/http"
	"strings"
)

type Registry interface {
	Name() string

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

func FindRegistry(repo string) Registry {
	for _, registry := range registries {
		if strings.HasPrefix(repo, registry.Name()) {
			return registry
		}
	}
	return nil
}
