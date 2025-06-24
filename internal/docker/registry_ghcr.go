package docker

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gabe565.com/transsmute/internal/config"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-github/v73/github"
)

type AuthMethod uint8

const (
	AuthNone AuthMethod = iota
	AuthToken
	AuthApp
)

func NewGhcr(conf config.GHCR) (*GHCR, error) {
	ghcr := &GHCR{
		username: conf.Username,
		password: conf.Password,

		installationID: conf.InstallationID,
	}
	appID := conf.AppID
	privateKey := []byte(conf.PrivateKey)
	if len(privateKey) == 0 {
		if conf.PrivateKeyPath != "" {
			var err error
			privateKey, err = os.ReadFile(conf.PrivateKeyPath)
			if err != nil {
				return ghcr, err
			}
		}
	}

	if ghcr.authMethod == AuthNone && appID != 0 && ghcr.installationID != 0 && len(privateKey) != 0 {
		itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, privateKey)
		if err != nil {
			return ghcr, err
		}

		ghcr.client = github.NewClient(&http.Client{Transport: itr})
		ghcr.username = strconv.Itoa(int(ghcr.installationID))
		ghcr.authMethod = AuthApp
		if err := ghcr.RefreshAppToken(context.Background()); err != nil {
			return ghcr, err
		}
	}

	if ghcr.authMethod == AuthNone && ghcr.username != "" && ghcr.password != "" {
		ghcr.authMethod = AuthToken
	}

	return ghcr, nil
}

type GHCR struct {
	authMethod AuthMethod

	username string
	password string

	client         *github.Client
	installationID int64
	expiresAt      time.Time
}

func (g *GHCR) Match(repo string) bool {
	return strings.HasPrefix(repo, "ghcr.io/")
}

func (g *GHCR) Authenticator(ctx context.Context, _ string) (authn.Authenticator, error) {
	if g.authMethod == AuthApp && time.Until(g.expiresAt) < time.Minute {
		if err := g.RefreshAppToken(ctx); err != nil {
			return nil, err
		}
	}

	return &authn.Basic{
		Username: g.username,
		Password: g.password,
	}, nil
}

func (g *GHCR) GetRepoURL(repo string) *url.URL {
	return &url.URL{Scheme: "https", Host: "ghcr.io", Path: strings.TrimPrefix(repo, "ghcr.io/")}
}

func (g *GHCR) GetTagURL(repo, _ string) *url.URL {
	return g.GetRepoURL(repo)
}

func (g *GHCR) GetOwner(repo string) string {
	return strings.SplitN(repo, "/", 3)[1]
}

func (g *GHCR) RefreshAppToken(ctx context.Context) error {
	token, _, err := g.client.Apps.CreateInstallationToken(ctx, g.installationID, &github.InstallationTokenOptions{
		Permissions: &github.InstallationPermissions{
			Packages: github.Ptr("read"),
		},
	})
	if err != nil {
		return err
	}

	g.password = token.GetToken()
	if token.ExpiresAt != nil {
		g.expiresAt = *token.ExpiresAt.GetTime()
	}
	return nil
}
