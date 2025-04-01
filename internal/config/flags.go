package config

import "github.com/spf13/cobra"

const (
	ListenAddressFlag = "address"

	DockerEnabledFlag      = "docker-enabled"
	GHCRUsernameFlag       = "ghcr-username"
	GHCRPasswordFlag       = "ghcr-password"
	GHCRAppIDFlag          = "ghcr-app-id"
	GHCRInstallationIDFlag = "ghcr-installation-id"
	GHCRPrivateKeyFlag     = "ghcr-private-key"
	GHCRPrivateKeyPathFlag = "ghcr-private-key-path"
	DockerHubUsernameFlag  = "dockerhub-username"
	DockerHubPasswordFlag  = "dockerhub-password" //nolint:gosec

	YouTubeEnabledFlag = "youtube-enabled"
	YouTubeAPIKeyFlag  = "youtube-key"

	KemonoEnabledFlag = "kemono-enabled"
	KemonoHostsFlag   = "kemono-hosts"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.ListenAddress, ListenAddressFlag, c.ListenAddress, "Listen address")

	cmd.Flags().BoolVar(&c.Docker.Enabled, DockerEnabledFlag, c.Docker.Enabled, "Docker API enabled")
	cmd.Flags().StringVar(&c.Docker.GHCR.Username, GHCRUsernameFlag, c.Docker.GHCR.Username,
		"GitHub username for ghcr.io",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.Password, GHCRPasswordFlag, c.Docker.GHCR.Password,
		"GitHub personal access token for ghcr.io",
	)
	cmd.Flags().Int64Var(&c.Docker.GHCR.AppID, GHCRAppIDFlag, c.Docker.GHCR.AppID, "GitHub app ID")
	cmd.Flags().Int64Var(&c.Docker.GHCR.InstallationID, GHCRInstallationIDFlag, c.Docker.GHCR.InstallationID,
		"GitHub installation ID",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.PrivateKey, GHCRPrivateKeyFlag, c.Docker.GHCR.PrivateKey,
		"GitHub app private key",
	)
	cmd.Flags().StringVar(&c.Docker.GHCR.PrivateKeyPath, GHCRPrivateKeyPathFlag, c.Docker.GHCR.PrivateKeyPath,
		"GitHub app private key file path",
	)
	cmd.Flags().StringVar(&c.Docker.DockerHub.Username, DockerHubUsernameFlag, c.Docker.DockerHub.Username,
		"DockerHub username for private repos",
	)
	cmd.Flags().StringVar(&c.Docker.DockerHub.Password, DockerHubPasswordFlag, c.Docker.DockerHub.Password,
		"DockerHub password for private repos",
	)

	cmd.Flags().BoolVar(&c.YouTube.Enabled, YouTubeEnabledFlag, c.YouTube.Enabled, "YouTube API enabled")
	cmd.Flags().StringVar(&c.YouTube.APIKey, YouTubeAPIKeyFlag, c.YouTube.APIKey, "YouTube API key")

	cmd.Flags().BoolVar(&c.Kemono.Enabled, KemonoEnabledFlag, c.Kemono.Enabled, "Kemono API enabled")
	cmd.Flags().StringToStringVar(&c.Kemono.Hosts, KemonoHostsFlag, c.Kemono.Hosts,
		"Kemono API hosts, where the key is the URL prefix and the value is the host",
	)
}
