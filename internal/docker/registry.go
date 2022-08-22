package docker

import (
	"net/http"
	"strings"
)

type Registry interface {
	Name() string

	ApiUrl() string
	TokenUrl(repo string) string

	Transport(repo string) http.RoundTripper

	NormalizeRepo(repo string) string
	GetRepoUrl(repo string) string
	GetTagUrl(repo, tag string) string
	GetOwner(repo string) string
}

var Registries = []Registry{
	Ghcr{},
	DockerHub{},
}

func FindRegistry(repo string) Registry {
	for _, registry := range Registries {
		if strings.HasPrefix(repo, registry.Name()) {
			return registry
		}
	}
	return nil
}
