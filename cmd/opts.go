package cmd

import (
	"github.com/spf13/cobra"
)

type Option func(cmd *cobra.Command)

func WithVersion(version string) Option {
	return func(cmd *cobra.Command) {
		if cmd.Annotations == nil {
			cmd.Annotations = make(map[string]string, 2)
		}
		cmd.Annotations["version"] = version
		cmd.Version, cmd.Annotations["commit"] = buildVersion(version)
	}
}
