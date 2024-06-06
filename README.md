# Transsmute

<img src="./assets/icon.svg" alt="Transsmute Icon" width="92" align="right">

[![Build](https://github.com/gabe565/transsmute/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/transsmute/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/transsmute)

Transsmute (transmute + RSS) is a server written in Go that builds RSS
feeds for websites that don't provide them.

Currently, the following feed types are supported:

- YouTube (channel, playlist)
- Container registries (DockerHub, ghcr.io)

## Installation

Transsmute can be installed as a container (suggested) or as a compiled
binary.

### Docker

A Docker container is available at `ghcr.io/gabe565/transsmute`. The
containerized version of Transsmute will run on port 80 by default,
and accepts all configuration as flags or as environment variables.
See [Configuration](#configuration) for more info.

```shell
docker run --rm -it -p 80:80 ghcr.io/gabe565/transsmute
```

Alternatively, an example [`docker-compose.yml`](/docker-compose.yml) file
is provided.

### Kubernetes

A Helm chart is available for Kubernetes deployment.
For more information, see
[charts.gabe565.com](https://charts.gabe565.com/charts/transsmute/) or
[Artifact Hub](https://artifacthub.io/packages/helm/gabe565/transsmute).

## Configuration

| Flag                   | Environment Variable            | Description                                                                  | Default                          |
|------------------------|---------------------------------|------------------------------------------------------------------------------|----------------------------------|
| `--address`            | `TRANSSMUTE_ADDRESS`            | Listen address                                                               | `":3000"` (`":80"` in container) |
| `--youtube-enabled`    | `TRANSSMUTE_YOUTUBE_ENABLED`    | YouTube API enabled.                                                         | `true`                           |
| `--youtube-key`        | `TRANSSMUTE_YOUTUBE_KEY`        | YouTube API key. Required to enable YouTube routes!                          | `""`                             |
| `--docker-enabled`     | `TRANSSMUTE_DOCKER_ENABLED`     | Docker API enabled.                                                          | `true`                           |
| `--dockerhub-username` | `TRANSSMUTE_DOCKERHUB_USERNAME` | DockerHub username for private repos.                                        | `""`                             |
| `--dockerhub-password` | `TRANSSMUTE_DOCKERHUB_PASSWORD` | DockerHub password for private repos.                                        | `""`                             |
| `--ghcr-username`      | `TRANSSMUTE_GHCR_USERNAME`      | GitHub username for [ghcr.io](https://ghcr.io) repos.                        | `""`                             |
| `--ghcr-password`      | `TRANSSMUTE_GHCR_PASSWORD`      | GitHub PAT for [ghcr.io](https://ghcr.io) repos.                             | `""`                             |
| `--kemono-enabled`     | `TRANSSMUTE_KEMONO_ENABLED`     | Kemono API enabled.                                                          | `true`                           |
| `--kemono-hosts`       | `TRANSSMUTE_KEMONO_HOSTS`       | Kemono API hosts, where the key is the URL prefix and the value is the host. | `kemono=kemono.su`               |

### DockerHub

DockerHub credentials are only required to access private repositories.

### ghcr.io

A personal access token is used to authenticate into GitHub's ghcr.io API.
The only required scope is `read:packages`.
[Click here](https://github.com/settings/tokens/new?description=Transsmute&scopes=read:packages)
to generate a personal access token with the necessary scopes prefilled.

## Routes

### Feed Type

An Atom feed is generated by default, but a file extension of
`.json` or `.rss` will change the output to the given format.

### YouTube

- `/youtube/playlist/:playlistId`
- `/youtube/channel/:channelId`

### Docker

- `/docker/tags/:repo`

### Kemono

- `/kemono/{service}/user/{name}`
