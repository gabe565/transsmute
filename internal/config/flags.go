package config

import "github.com/spf13/cobra"

const (
	FlagListenAddress = "address"
	FlagTLSCertPath   = "tls-cert-path"
	FlagTLSKeyPath    = "tls-key-path"

	FlagDockerEnabled      = "docker-enabled"
	FlagGHCRUsername       = "ghcr-username"
	FlagGHCRPassword       = "ghcr-password"
	FlagGHCRAppID          = "ghcr-app-id"
	FlagGHCRInstallationID = "ghcr-installation-id"
	FlagGHCRPrivateKey     = "ghcr-private-key"
	FlagGHCRPrivateKeyPath = "ghcr-private-key-path"
	FlagDockerHubUsername  = "dockerhub-username"
	FlagDockerHubPassword  = "dockerhub-password" //nolint:gosec

	FlagYouTubeEnabled = "youtube-enabled"
	FlagYouTubeAPIKey  = "youtube-key"

	FlagKemonoEnabled = "kemono-enabled"
	FlagKemonoHosts   = "kemono-hosts"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.ListenAddress, FlagListenAddress, c.ListenAddress, "Listen address")
	cmd.Flags().StringVar(&c.TLSCertPath, FlagTLSCertPath, c.TLSCertPath, "TLS certificate path for HTTPS listener")
	cmd.Flags().StringVar(&c.TLSKeyPath, FlagTLSKeyPath, c.TLSKeyPath, "TLS key path for HTTPS listener")

	cmd.Flags().BoolVar(&c.Docker.Enabled, FlagDockerEnabled, c.Docker.Enabled, "Docker API enabled")
	cmd.Flags().StringVar(&c.Docker.GHCR.Username, FlagGHCRUsername, c.Docker.GHCR.Username,
		"GitHub username for ghcr.io",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.Password, FlagGHCRPassword, c.Docker.GHCR.Password,
		"GitHub personal access token for ghcr.io",
	)
	cmd.Flags().Int64Var(&c.Docker.GHCR.AppID, FlagGHCRAppID, c.Docker.GHCR.AppID, "GitHub app ID")
	cmd.Flags().Int64Var(&c.Docker.GHCR.InstallationID, FlagGHCRInstallationID, c.Docker.GHCR.InstallationID,
		"GitHub installation ID",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.PrivateKey, FlagGHCRPrivateKey, c.Docker.GHCR.PrivateKey,
		"GitHub app private key",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.PrivateKeyPath, FlagGHCRPrivateKeyPath, c.Docker.GHCR.PrivateKeyPath,
		"GitHub app private key file path",
	)
	cmd.Flags().StringVar(&c.Docker.DockerHub.Username, FlagDockerHubUsername, c.Docker.DockerHub.Username,
		"DockerHub username for private repos",
	)
	cmd.Flags().StringVar(&c.Docker.DockerHub.Password, FlagDockerHubPassword, c.Docker.DockerHub.Password,
		"DockerHub password for private repos",
	)

	cmd.Flags().BoolVar(&c.YouTube.Enabled, FlagYouTubeEnabled, c.YouTube.Enabled, "YouTube API enabled")
	cmd.Flags().StringVar(&c.YouTube.APIKey, FlagYouTubeAPIKey, c.YouTube.APIKey, "YouTube API key")

	cmd.Flags().BoolVar(&c.Kemono.Enabled, FlagKemonoEnabled, c.Kemono.Enabled, "Kemono API enabled")
	cmd.Flags().StringToStringVar(&c.Kemono.Hosts, FlagKemonoHosts, c.Kemono.Hosts,
		"Kemono API hosts, where the key is the URL prefix and the value is the host",
	)
}
