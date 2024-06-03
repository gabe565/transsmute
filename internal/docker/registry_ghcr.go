package docker

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-github/v62/github"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AuthMethod uint8

const (
	AuthNone AuthMethod = iota
	AuthToken
	AuthApp
)

func init() {
	flag.String("ghcr-username", "", "GitHub username for ghcr.io")
	if err := viper.BindPFlag("ghcr.username", flag.Lookup("ghcr-username")); err != nil {
		panic(err)
	}

	flag.String("ghcr-password", "", "GitHub personal access token for ghcr.io")
	if err := viper.BindPFlag("ghcr.password", flag.Lookup("ghcr-password")); err != nil {
		panic(err)
	}

	flag.Int64("ghcr-app-id", 0, "GitHub app ID")
	if err := viper.BindPFlag("ghcr.app-id", flag.Lookup("ghcr-app-id")); err != nil {
		panic(err)
	}

	flag.Int64("ghcr-installation-id", 0, "GitHub installation ID")
	if err := viper.BindPFlag("ghcr.installation-id", flag.Lookup("ghcr-installation-id")); err != nil {
		panic(err)
	}

	flag.String("ghcr-private-key", "", "GitHub app private key")
	if err := viper.BindPFlag("ghcr.private-key", flag.Lookup("ghcr-private-key")); err != nil {
		panic(err)
	}

	flag.String("ghcr-private-key-path", "", "GitHub app private key file path")
	if err := viper.BindPFlag("ghcr.private-key-path", flag.Lookup("ghcr-private-key-path")); err != nil {
		panic(err)
	}
}

func NewGhcr() (*Ghcr, error) {
	ghcr := &Ghcr{
		username: viper.GetString("ghcr.username"),
		password: viper.GetString("ghcr.password"),

		installationId: viper.GetInt64("ghcr.installation-id"),
	}
	appId := viper.GetInt64("ghcr.app-id")
	privateKey := []byte(viper.GetString("ghcr.private-key"))
	if len(privateKey) == 0 {
		privateKeyPath := viper.GetString("ghcr.private-key-path")
		if privateKeyPath != "" {
			var err error
			privateKey, err = os.ReadFile(privateKeyPath)
			if err != nil {
				return ghcr, err
			}
		}
	}

	if ghcr.authMethod == AuthNone && appId != 0 && ghcr.installationId != 0 && len(privateKey) != 0 {
		itr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appId, privateKey)
		if err != nil {
			return ghcr, err
		}

		ghcr.client = github.NewClient(&http.Client{Transport: itr})
		ghcr.username = strconv.Itoa(int(ghcr.installationId))
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

type Ghcr struct {
	authMethod AuthMethod

	username string
	password string

	client         *github.Client
	installationId int64
	expiresAt      time.Time
}

func (g Ghcr) Names() []string {
	return []string{"ghcr.io"}
}

func (g Ghcr) ApiUrl() string {
	return "https://ghcr.io"
}

func (g Ghcr) TokenUrl(repo string) string {
	return g.ApiUrl() + "/token?service=ghcr.io&scope=repository:" + repo + ":pull"
}

func (g Ghcr) Authenticator(ctx context.Context, _ string) (authn.Authenticator, error) {
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

func (g Ghcr) NormalizeRepo(repo string) string {
	return repo
}

func (g Ghcr) GetRepoUrl(repo string) string {
	return "https://" + repo
}

func (g Ghcr) GetTagUrl(repo, tag string) string {
	return g.GetRepoUrl(repo)
}

func (g Ghcr) GetOwner(repo string) string {
	return strings.SplitN(repo, "/", 3)[1]
}

func (g *Ghcr) RefreshAppToken(ctx context.Context) error {
	token, _, err := g.client.Apps.CreateInstallationToken(ctx, g.installationId, &github.InstallationTokenOptions{
		Permissions: &github.InstallationPermissions{
			Packages: github.String("read"),
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
